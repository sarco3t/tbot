/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	geoestclient "github.com/sarco3t/kbot/geoest_client"
	"github.com/spf13/cobra"
	"gopkg.in/telebot.v4"
)

var GEOEST_TEMPLATE = template.Must(template.New("response").Parse(`Prediction:
  Latitude: {{printf "%.4f" .Prediction.Latitude}}
  Longitude: {{printf "%.4f" .Prediction.Longitude}}
  Confidence: {{printf "%.2f" .Confidence}}%`))

var (
	// TeleToken bot
	TeleToken = os.Getenv("TELE_TOKEN")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Start the Kbot Telegram bot",
	Long:    "Kbot is a Telegram bot for geo estimation using photos.",
	Run: func(cmd *cobra.Command, args []string) {
		if TeleToken == "" {
			log.Fatal("TELE_TOKEN environment variable is not set")
		}

		fmt.Printf("kbot %s started\n", appVersion)

		// Initialize and start the bot
		if err := startBot(); err != nil {
			log.Fatalf("Failed to start bot: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}

// startBot initializes and starts the Telegram bot
func startBot() error {
	kbot, err := telebot.NewBot(telebot.Settings{
		URL:    "",
		Token:  TeleToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	// Handle text messages
	kbot.Handle(telebot.OnText, handleText)

	// Handle photo messages
	kbot.Handle(telebot.OnPhoto, handlePhoto)

	kbot.Start()
	return nil
}

// handleText processes text messages
func handleText(m telebot.Context) error {
	log.Printf("Received message: %s", m.Text())
	payload := m.Message().Payload
	switch payload {
	case "hello":
		return m.Send(fmt.Sprintf("Hello, I'm Kbot %s!", appVersion))
	default:
		return m.Send("Usage: Send a photo for geo estimation.\n\n")
	}
}

// handlePhoto processes photo messages
func handlePhoto(m telebot.Context) error {
	photo := m.Message().Photo

	// Create a geoestimation client
	client := geoestclient.NewClient("http://localhost:8000")

	reader, err := m.Bot().File(&photo.File)
	if err != nil {
		log.Printf("Error getting file: %s", err)
		return m.Reply("Failed to retrieve the photo. Please try again.")
	}
	// Evaluate the photo
	response, err := client.Evaluate(reader)
	if err != nil {
		log.Printf("Error evaluating photo: %s", err)
		return m.Reply("Failed to process the photo. Please try again.")
	}

	// Format the response using the template
	responseText, err := formatResponse(response)
	if err != nil {
		log.Printf("Error formatting response: %s", err)
		return m.Reply("Failed to format the response. Please try again.")
	}

	// Send the response
	return m.Reply(responseText)
}

// formatResponse formats the geo estimation response using a template
func formatResponse(response *geoestclient.UploadResponse) (string, error) {
	var responseText bytes.Buffer
	err := GEOEST_TEMPLATE.Execute(&responseText, response)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return responseText.String(), nil
}
