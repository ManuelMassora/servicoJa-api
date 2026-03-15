# servicoJa API

> **Micro SaaS backend** for service provision management — connecting clients to service providers through a clean, scalable Go API.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Running Tests](#running-tests)
- [Roadmap](#roadmap)
- [Author](#author)

---

## Overview

**servicoJa** is a micro SaaS platform focused on deepening the service provision experience. The API serves as the backbone of the platform, handling service listings, provider management, client requests, and bookings.

Built with Go following **Clean Architecture** principles, the project is designed for high performance, maintainability, and easy horizontal scaling.

Key design goals:
- Low-latency RESTful API
- Clear separation of concerns via `internal/` and `pkg/` layers
- Containerized and ready for CI/CD deployment
- End-to-end tested from day one

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go (Golang) |
| Architecture | Clean Architecture |
| Containerization | Docker |
| Testing | Go testing + E2E (`test/e2e`) |
| Package management | Go Modules (`go.mod`) |

---

## Project Structure

```
servicoJa-api/
├── internal/           # Private application code (domain, use cases, handlers)
├── pkg/                # Reusable public packages (shared utilities, middleware, etc.)
├── test/
│   └── e2e/            # End-to-end tests
├── main.go             # Application entry point
├── Dockerfile          # Container definition
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── TODO.md             # Planned features and improvements
└── .gitignore
```

The `internal/` directory follows standard Go project conventions — all core business logic (domain entities, use cases, HTTP handlers, repositories) lives here and is intentionally unexported. The `pkg/` directory holds reusable utilities that could theoretically be shared across services.

---

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (for containerized setup)
- [Git](https://git-scm.com/)

### Running with Docker

```bash
# Clone the repository
git clone https://github.com/ManuelMassora/servicoJa-api.git
cd servicoJa-api

# Build and run the container
docker build -t servicoja-api .
docker run -p 8080:8080 --env-file .env servicoja-api
```

The API will be available at `http://localhost:8080`.

### Running Locally

```bash
# Clone the repository
git clone https://github.com/ManuelMassora/servicoJa-api.git
cd servicoJa-api

# Install dependencies
go mod download

# Run the application
go run main.go
```

---

## Environment Variables

Create a `.env` file in the root of the project. Below is a reference for the expected variables:

```env
# Server
PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=servicoja

# Auth
JWT_SECRET=your_jwt_secret
```

> Never commit your `.env` file. It is already listed in `.gitignore`.

---

## API Endpoints

> Full API documentation coming soon. Below is a high-level overview of the planned resource structure.

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | API health check |
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Authenticate and receive a token |
| `GET` | `/api/v1/services` | List available services |
| `POST` | `/api/v1/services` | Create a new service listing |
| `GET` | `/api/v1/services/:id` | Get service details |
| `POST` | `/api/v1/bookings` | Request a service booking |
| `GET` | `/api/v1/bookings/:id` | Get booking status |

---

## Running Tests

### Unit & Integration Tests

```bash
go test ./...
```

### End-to-End Tests

```bash
go test ./test/e2e/...
```

To run with verbose output:

```bash
go test -v ./...
```

---

## Roadmap

See [TODO.md](./TODO.md) for the full list of planned features. High-level items include:

- [ ] Complete authentication flow (JWT)
- [ ] Service provider profile management
- [ ] Booking lifecycle (request → accept → complete → review)
- [ ] Notifications system
- [ ] API documentation (Swagger/OpenAPI)
- [ ] CI/CD pipeline configuration

---

## Author

**Manuel Massora** — Backend Engineer  
Maputo, Mozambique

- GitHub: [@ManuelMassora](https://github.com/ManuelMassora)
- LinkedIn: [manuelt-massora-5bb417375](https://linkedin.com/in/manuelt-massora-5bb417375/)
- Email: manuelmassora75@gmail.com

---

> Built with Go · Clean Architecture · Docker
