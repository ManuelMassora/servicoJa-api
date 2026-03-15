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



#####################


# servicoJa API

> **Backend Micro SaaS** para gestão de prestação de serviços — conectando clientes a prestadores através de uma API Go limpa e escalável.

---

## Índice

- [Visão Geral](#visão-geral)
- [Stack Tecnológica](#stack-tecnológica)
- [Estrutura do Projecto](#estrutura-do-projecto)
- [Como Começar](#como-começar)
  - [Pré-requisitos](#pré-requisitos)
  - [Executar com Docker](#executar-com-docker)
  - [Executar Localmente](#executar-localmente)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Endpoints da API](#endpoints-da-api)
- [Executar Testes](#executar-testes)
- [Roadmap](#roadmap)
- [Autor](#autor)

---

## Visão Geral

**servicoJa** é uma plataforma Micro SaaS focada em aprofundar a experiência de prestação de serviços. A API é a espinha dorsal da plataforma, responsável por gerir listagens de serviços, perfis de prestadores, pedidos de clientes e reservas.

Construída em Go seguindo os princípios de **Clean Architecture**, o projecto é desenhado para alto desempenho, facilidade de manutenção e escalabilidade horizontal.

Objectivos de design:
- API RESTful de baixa latência
- Separação clara de responsabilidades via camadas `internal/` e `pkg/`
- Contenorizado e pronto para pipelines de CI/CD
- Testado de ponta a ponta desde o início do desenvolvimento

---

## Stack Tecnológica

| Camada | Tecnologia |
|---|---|
| Linguagem | Go (Golang) |
| Arquitectura | Clean Architecture |
| Contenorização | Docker |
| Testes | Go testing + E2E (`test/e2e`) |
| Gestão de pacotes | Go Modules (`go.mod`) |

---

## Estrutura do Projecto

```
servicoJa-api/
├── internal/           # Código privado da aplicação (domínio, casos de uso, handlers)
├── pkg/                # Pacotes reutilizáveis (utilitários partilhados, middleware, etc.)
├── test/
│   └── e2e/            # Testes de ponta a ponta (end-to-end)
├── main.go             # Ponto de entrada da aplicação
├── Dockerfile          # Definição do contentor
├── go.mod              # Dependências do módulo Go
├── go.sum              # Checksums das dependências
├── TODO.md             # Funcionalidades planeadas e melhorias
└── .gitignore
```

O directório `internal/` segue as convenções padrão de projectos Go — toda a lógica de negócio principal (entidades de domínio, casos de uso, handlers HTTP, repositórios) reside aqui e é intencionalmente não exportada. O directório `pkg/` contém utilitários reutilizáveis que poderiam ser partilhados entre serviços.

---

## Como Começar

### Pré-requisitos

Certifica-te de que tens instalado:

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (para configuração contenorizada)
- [Git](https://git-scm.com/)

### Executar com Docker

```bash
# Clonar o repositório
git clone https://github.com/ManuelMassora/servicoJa-api.git
cd servicoJa-api

# Construir e executar o contentor
docker build -t servicoja-api .
docker run -p 8080:8080 --env-file .env servicoja-api
```

A API ficará disponível em `http://localhost:8080`.

### Executar Localmente

```bash
# Clonar o repositório
git clone https://github.com/ManuelMassora/servicoJa-api.git
cd servicoJa-api

# Instalar dependências
go mod download

# Executar a aplicação
go run main.go
```

---

## Variáveis de Ambiente

Cria um ficheiro `.env` na raiz do projecto. Abaixo está uma referência das variáveis esperadas:

```env
# Servidor
PORT=8080
APP_ENV=development

# Base de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=o_teu_utilizador
DB_PASSWORD=a_tua_password
DB_NAME=servicoja

# Autenticação
JWT_SECRET=o_teu_segredo_jwt
```

> Nunca faças commit do ficheiro `.env`. Já está listado no `.gitignore`.

---

## Endpoints da API

> Documentação completa da API em breve. Abaixo está uma visão geral de alto nível da estrutura de recursos planeada.

| Método | Endpoint | Descrição |
|---|---|---|
| `GET` | `/health` | Verificação de saúde da API |
| `POST` | `/api/v1/auth/register` | Registar um novo utilizador |
| `POST` | `/api/v1/auth/login` | Autenticar e receber token |
| `GET` | `/api/v1/services` | Listar serviços disponíveis |
| `POST` | `/api/v1/services` | Criar uma nova listagem de serviço |
| `GET` | `/api/v1/services/:id` | Obter detalhes de um serviço |
| `POST` | `/api/v1/bookings` | Solicitar uma reserva de serviço |
| `GET` | `/api/v1/bookings/:id` | Obter estado de uma reserva |

---

## Executar Testes

### Testes Unitários e de Integração

```bash
go test ./...
```

### Testes End-to-End

```bash
go test ./test/e2e/...
```

Para executar com saída detalhada:

```bash
go test -v ./...
```

---

## Roadmap

Consulta o [TODO.md](./TODO.md) para a lista completa de funcionalidades planeadas. Itens principais:

- [ ] Fluxo de autenticação completo (JWT)
- [ ] Gestão de perfil do prestador de serviços
- [ ] Ciclo de vida de reservas (pedido → aceitação → conclusão → avaliação)
- [ ] Sistema de notificações
- [ ] Documentação da API (Swagger/OpenAPI)
- [ ] Configuração de pipeline CI/CD

---

## Autor

**Manuel Massora** — Backend Engineer  
Maputo, Moçambique

- GitHub: [@ManuelMassora](https://github.com/ManuelMassora)
- LinkedIn: [manuelt-massora-5bb417375](https://linkedin.com/in/manuelt-massora-5bb417375/)
- Email: manuelmassora75@gmail.com

---

> Construído com Go · Clean Architecture · Docker


