# Subscriptions Aggregator API

Test assignment for Golang Developer position.

REST service for managing user subscriptions and calculating total cost for a selected period.

---

## Features

- CRUD operations for subscriptions
- Total cost calculation with filters
- PostgreSQL storage
- Database migrations
- Swagger API documentation
- Docker Compose setup
- Graceful shutdown
- Structured logging

---

## Tech Stack

- Go 1.25
- Gin
- PostgreSQL
- golang-migrate
- Swagger (swaggo)
- Docker / Docker Compose

---

## Running the project

### Requirements
- Docker
- Docker Compose

### Start service

### bash
docker compose up --build

### Service will be available at:
http://localhost:8080

### Swagger API:
http://localhost:8080/swagger/index.html


## Configuration

### Configuration is loaded from config.yaml.
Environment variables can be provided via .env file.
Example configuration is available in .env.example.


## Architecture
### Project follows a layered architecture:
- handlers — HTTP layer (Gin handlers, DTOs)
- service — business logic
- repo — database access (PostgreSQL)
- domain — core entities

This separation allows easy testing and maintenance.

## Database
PostgreSQL is started via Docker Compose.
Database schema is initialized using migrations on startup.

## Health check
GET /health
Returns service and database status.