# Docker Setup and Running Instructions

This project is containerized using Docker and Docker Compose. The stack consists of three components:

1. **Backend**: A Go REST API server that binds to port 8080
2. **Pulsesim**: A Go client that simulates online activity for 100 users
3. **Dashboard**: A Vite + Svelte webapp that shows online status for the first 100 users

## Prerequisites

- Docker
- Docker Compose

## Running the Stack

To build and run the entire stack:

```bash
docker-compose up --build
```

By default, the Go applications (backend and pulsesim) are built for Linux with amd64 architecture (x64 machines). You can specify a different target architecture by setting the GOARCH environment variable:

```bash
GOARCH=arm64 docker-compose up --build  # Build for ARM64 architecture (e.g., Apple M1/M2 laptops)
```

Supported GOARCH values include:
- amd64 (default, 64-bit x86)
- arm64 (64-bit ARM)
- 386 (32-bit x86)
- arm (32-bit ARM)

This will:
- Build and start the backend service on port 8080
- Build and start the pulsesim service, which will register and simulate 100 users
- Build and start the dashboard service on port 9001

You can access the dashboard at http://localhost:9001

## Services

- **Backend**: http://localhost:8080
- **Dashboard**: http://localhost:9001

## Stopping the Stack

To stop the stack:

```bash
docker-compose down
```
