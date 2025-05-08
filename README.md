# tbot

> **Telegram + AI visual-geo assistant**  
> Send a photo, get back an estimated location.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.111+-009688?logo=fastapi&logoColor=white)](https://fastapi.tiangolo.com/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Chat](https://img.shields.io/badge/telegram-@sasdasds1239999__bot-2CA5E0?logo=telegram&logoColor=white)](https://t.me/sasdasds1239999_bot)

## 🏃‍♂️ Quick start (local dev)
## ✨ Features

* Telegram bot written in **Go** (`cmd/`, `main.go`)
* **Visual-location prediction** via the `visloc-estimation` submodule (FastAPI service)
<!-- * Multi-stage Docker build for both Go and Python stacks (see `build/` & `Dockerfile`s) -->
* `.env`-driven configuration — no hard-coded secrets
* Makefile helpers for common tasks
<!-- * CI-ready: images are fully reproducible and run either locally (Docker Compose) or in production (Kubernetes/Fargate). -->
* (https://t.me/sasdasds1239999_bot)bot link
---

## 🗺️ Architecture

```text
Telegram ⇄ tbot (Go)
              │ HTTP/JSON
              │
              ├──▶ visloc-estimation (FastAPI, Python)
              │       └──– pre-trained CNN + ViT models
              │
              └──▶ GeoEst API / model client
```

---

## 🏃‍♂️ Quick start (local dev, no Docker)

### 1. Clone (with submodules)

```bash
git clone --recurse-submodules https://github.com/sarco3t/tbot.git
cd tbot

```

### 2. Pull dependencies

#### Go modules
```bash
go mod download
```
#### Python requirements
```bash
cd visloc-estimation
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
cd ..

```


### 3. Configure environment

1. Copy the template  
   ```bash
   cp .env.example .env
2. Open .env and fill in at least


### Run services
1. Terminal 1 – FastAPI (Python)

```bash
cd visloc-estimation
uvicorn api.server:app --host 0.0.0.0 --port 8000 # or ./entrypoint.sh
```
2.Terminal 2 – Telegram bot (Go)

```bash
make run 
```
Navigate back to [Telegram](https://t.me/sasdasds1239999_bot), send /start or a photo, and enjoy!

