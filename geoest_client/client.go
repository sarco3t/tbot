package geoestclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// Client represents the API client for interacting with the third-party service.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// UploadFile uploads a file to the third-party API.
func (c *Client) uploadFile(endpoint, fieldName string, file io.Reader) ([]byte, error) {
	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file to form: %w", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/%s", c.BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, nil
}

func (c *Client) Evaluate(file io.Reader) (*UploadResponse, error) {

	// Call the generic UploadFile method
	response, err := c.uploadFile("evaluate", "image", file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload photo: %w", err)
	}
	var uploadResponse UploadResponse
	err = json.Unmarshal(response, &uploadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &uploadResponse, nil
}

type UploadResponse struct {
	Prediction struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"prediction"`
	Confidence float64 `json:"confidence"`
}
