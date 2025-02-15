# Go Web Service with Kafka

This is a simple Go web service that integrates with Kafka. The project uses `air` for live-reloading of Go code during development and Docker Compose to set up Kafka.

## Prerequisites

Before you get started, make sure you have the following installed:

- [Go](https://golang.org/dl/) (1.18+)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- [air](https://github.com/air-verse/air) (for live-reloading Go code)
- **Operating System**: Linux or WSL (Windows Subsystem for Linux) on Windows.

If you are using Windows, it's recommended to use WSL for a seamless development experience.

## Getting Started

Follow the steps below to set up and run the project:

### 1. Initialize Kafka using Docker Compose

Start the Kafka services using the provided `docker-compose` file:

```bash
docker compose -f docker-compose-kafka.yaml up
```

This command will start the Kafka and Zookeeper containers required for the application.

### 2. Start the Server

Use the `air` command to build and run the server:

```bash
air
```