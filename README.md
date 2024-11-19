# Project Commands

## Overview

This project uses a Makefile to provide easy-to-use commands for development, testing, and deployment.

## Prerequisites

- Docker
- Docker Compose
- Go
- golangci-lint
- golang-migrate

## Available Commands

### Development

| Command | Description |
|---------|-------------|
| `make help` | Display all available commands |
| `make compose-up` | Start PostgreSQL and application containers |
| `make compose-down` | Stop and remove all containers |

### Testing

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests with Ginkgo |
| `make lint` | Run golangci-lint code quality checks |
| `make compose-up-integration-test` | Run integration tests |

### Database Migrations

| Command | Description | Example |
|---------|-------------|---------|
| `make migrate-create NAME=migration_name` | Create a new SQL migration | `make migrate-create NAME=add_users_table` |

## Quick Start

1. Clone the repository
2. Run `make compose-up` to start the application

## Running Tests

- Unit Tests: `make test`
- Integration Tests: `make compose-up-integration-test`
- Linting: `make lint`

## Creating Migrations

To create a new database migration:

```bash
make migrate-create NAME=your_migration_description