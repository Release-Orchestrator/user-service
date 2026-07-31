# User Service

A RESTful microservice for user management, part of the Release Orchestrator platform.

## Tech Stack

- **Language:** Go 1.22
- **Framework:** Gin
- **Database:** PostgreSQL (via pgx)
- **Config:** Environment variables + .env

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| POST | /api/v1/users | Create user |
| GET | /api/v1/users | List users |
| GET | /api/v1/users/:id | Get user |
| PUT | /api/v1/users/:id | Update user |
| DELETE | /api/v1/users/:id | Delete user |
| GET | /internal/users/:id | Validate user exists (returns `{"exists": bool}`) |

## Quick Start

### Prerequisites

- Go 1.22+
- Docker (for PostgreSQL)

### Run locally

```bash
# Start PostgreSQL
docker compose up -d

# Run the service
make run
```

### Build Docker image

```bash
make docker-build VERSION=0.1.0
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | HTTP port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASS | postgres | Database password |
| DB_NAME | user_db | Database name |

## Project Structure

```
.
├── cmd/main.go              # Entry point
├── internal/
│   ├── config/              # Configuration
│   ├── handler/             # HTTP handlers
│   ├── model/               # Data models
│   ├── repository/          # Database layer
│   └── service/             # Business logic
├── migrations/              # SQL migrations
├── helm/                    # Kubernetes deployment
└── Dockerfile               # Multi-stage build
```
