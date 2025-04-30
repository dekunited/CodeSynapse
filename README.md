# CodeSynapse
by Manan Patel, Zachary Perry, Shayana Shreshta, Eric Vaughan

## Overview

This project looks to explore the capabilities of LLMs in cross-language code translation and its accuracy when translating across languages with different paradigms. In addition, we aim to build a tool that will allow developers to simply input code and translate it to a desired language using the LLM that is best at that specific translation.


## 🚀 Getting Started


### 📋 Prerequisites
Before you begin, make sure you have the following installed:
- **[Docker](https://www.docker.com/products/docker-desktop)**
  - Needed to run the containers for both the frontend + backend of the web app
- **[Go](https://golang.org/dl/)**
  - Required for building and running the backend code.
- **[Node.js](https://nodejs.org/en/download/)**
  - Required for building and running the frontend code.
- **[pnpm](https://pnpm.io/)**
  - Install with: `npm install -g pnpm` (or `brew install pnpm` for macOS).

### ⚙️ Development Workflow + How to Run
To run the containers, you can utilize the provided shell script: `./run_dev.sh`
- This will build both the frontend + backend containers with `docker compose build dir_name`
- Then, it will run everything with `docker compose up`
- When interrupted with `ctrl+c`, it will tear down everything with `docker compose down`
- Please note that both the frontend and backend have hot reload enabled (vite, air), so there's no need to teardown everything when adding changes. Only time you will need to rebuild the container is if you install any new dependencies. 
