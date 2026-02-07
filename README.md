# Web Forum (Go)

📄 [Русская версия README](README.ru.md)

A feature-rich web forum application built with Go.  
The project demonstrates clean layered architecture, server-side rendering, authentication, moderation workflows, and secure HTTP handling using the Go standard library with minimal external dependencies.

---

## What This Project Demonstrates

This project showcases practical Go backend development and architectural decisions commonly used in production systems.

It demonstrates:

- Layered architecture with clear separation of concerns (handler, service, repository)
- Clean dependency injection without frameworks
- Server-side rendering using Go `html/template`
- Authentication and authorization with session management
- OAuth 2.0 integration (Google and GitHub)
- Moderation and approval workflows
- SQLite database design with foreign key constraints
- HTTPS setup with TLS
- Dockerized deployment workflow

---

## Core Features

- User registration and authentication
- OAuth login via Google and GitHub
- Role system: user, moderator, administrator
- Post creation with categories and optional images
- Threaded comments
- Like and dislike reactions
- Post approval workflow
- Content reporting and moderation
- Notification system
- Secure sessions and CSRF protection

---

## Architecture Overview

The application follows a layered architecture pattern:

- **Handler layer**: HTTP request handling, routing, and response rendering
- **Service layer**: business logic and use cases
- **Repository layer**: database access abstraction
- **Entities**: domain models
- **UI layer**: embedded HTML templates and static assets

Each layer communicates through interfaces, allowing independent testing and future replacement of implementations.

---

## Project Structure

```
cmd/          – application entry point
internal/     – core application logic
  handler/    – HTTP handlers and routing
  service/    – business logic
  repository/ – data access
  entities/   – domain models
  session/    – session management
pkg/          – reusable packages and configuration
ui/           – HTML templates and static assets
schema/       – database schema
```

This structure follows Go conventions and keeps implementation details encapsulated inside the `internal` package.

---

## Technology Stack

- Go 1.22
- SQLite
- `net/http`
- `html/template`
- OAuth 2.0
- Docker
- TLS (HTTPS)

---

## Running Locally

### Requirements

- Go 1.22+
- Git
- SQLite
- OpenSSL

### Setup

Clone the repository:

```bash
git clone https://github.com/asari92/forum
cd forum
```

Initialize database and TLS certificates:

```bash
make init
```

Run the server:

```bash
go run ./cmd/web/main.go
```

The application will be available at:

https://localhost:4000

---

## Running with Docker

Build the image:

```bash
make build
```

Run the container:

```bash
make run
```

Stop the container:

```bash
make stop
```

---

## Notes

Technology choices were made intentionally to emphasize core Go concepts: explicit dependency management, clear separation of concerns, and predictable HTTP and data-access flows.
