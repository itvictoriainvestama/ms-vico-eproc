# E-Procurement API

Backend API for the E-Procurement system used to manage procurement workflows such as Purchase Request (PR), Request for Quotation (RFQ), Purchase Order (PO), vendor management, approval tasks, and authentication.

This project is written in Go, uses Gin as the HTTP framework, GORM as the ORM, and MySQL as the database.

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Current Modules](#current-modules)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Local Environment Variables](#local-environment-variables)
- [Setup Guide](#setup-guide)
- [Run the Project](#run-the-project)
- [Database Migration](#database-migration)
- [Build and Test](#build-and-test)
- [API Base URL](#api-base-url)
- [Available Routes](#available-routes)
- [Authentication](#authentication)
- [Database Notes](#database-notes)
- [Troubleshooting](#troubleshooting)
- [Development Notes](#development-notes)

## Overview

The API currently provides foundational procurement endpoints for:

- authentication
- purchase requisitions
- RFQ
- purchase orders
- vendor management
- approval tasks

The business reference for this project is documented in [docs/FSD_E-Procurement_Complete_2.docx.md](/Users/itvico/Dev/e-proc-api/docs/FSD_E-Procurement_Complete_2.docx.md).

## Tech Stack

- Go `1.22+` from [go.mod](/Users/itvico/Dev/e-proc-api/go.mod)
- Gin `github.com/gin-gonic/gin`
- GORM `gorm.io/gorm`
- MySQL `gorm.io/driver/mysql`
- JWT `github.com/golang-jwt/jwt/v5`
- `.env` loading with `github.com/joho/godotenv`

## Current Modules

Current implemented backend areas in this repository:

- `Auth`
- `PR`
- `RFQ`
- `PO`
- `Vendor`
- `Approval`

Important note:

- The FSD describes a broader target scope than what is currently implemented in code.
- Some advanced modules from the FSD, such as dynamic procurement policy, budget management, entity management, delegate approver, reference price, and vendor blacklist, are not fully implemented yet in this codebase.

## Project Structure

```text
e-proc-api/
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── BRD_E-Procurement.md
│   ├── FSD_E-Procurement_Complete_2.docx.md
│   └── SETUP_LOCAL.md
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── router/
│   └── services/
├── .env.example
├── Makefile
└── go.mod
```

## Prerequisites

Before running the project, make sure the following are installed on your machine:

- Go `1.22` or newer
- MySQL `8.x` or newer

Recommended local versions already verified for this project:

- Go `1.26.1`
- MySQL `9.6.0`

### Verify Installation

```bash
go version
mysql --version
```

If you use Homebrew on macOS:

```bash
brew install go mysql
brew services start mysql
```

## Local Environment Variables

The application loads configuration automatically from a `.env` file in the project root.

Start from [.env.example](/Users/itvico/Dev/e-proc-api/.env.example).

Example local `.env`:

```env
# Application
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=e_procurement

# Run auto migration on startup (set to true only in dev/first run)
DB_MIGRATE=false

# JWT
JWT_SECRET=local-dev-secret-change-me
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168
```

### Environment Variable Reference

| Variable | Description | Default |
| --- | --- | --- |
| `APP_PORT` | HTTP server port | `8080` |
| `APP_ENV` | Application environment | `development` |
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASSWORD` | MySQL password | empty |
| `DB_NAME` | MySQL database name | `e_procurement` |
| `DB_MIGRATE` | Run auto migration at startup | `false` |
| `DB_RESET` | Drop and recreate the configured database on startup | `false` |
| `DB_SEED` | Seed baseline entities, roles, departments, admin user, and user roles | `false` |
| `SEED_ADMIN_PASSWORD` | Password for the seeded admin user | `Admin123!` |
| `SEED_ENTITY_CODE` | Default seeded entity code | `HO` |
| `SEED_ENTITY_NAME` | Default seeded entity name | `Head Office` |
| `SEED_DEPARTMENT_CODE` | Default seeded department code | `PROC` |
| `SEED_DEPARTMENT_NAME` | Default seeded department name | `Procurement` |
| `JWT_SECRET` | JWT signing secret | `change-me-in-production` |
| `JWT_EXPIRY_HOURS` | Access token expiry in hours | `24` |
| `JWT_REFRESH_EXPIRY_HOURS` | Refresh token expiry in hours | `168` |

## Setup Guide

Follow these steps when setting up the project for the first time.

### 1. Clone the repository

```bash
git clone <repository-url>
cd e-proc-api
```

If you already have the source folder, just change into the project directory.

### 2. Create your environment file

Copy the example file:

```bash
cp .env.example .env
```

Then update `.env` to match your MySQL credentials.

### 3. Start MySQL

If you use Homebrew:

```bash
brew services start mysql
```

To confirm MySQL is running:

```bash
brew services list
mysql -u root -e "SELECT VERSION();"
```

### 4. Create the database

Create the database used by the application:

```bash
mysql -u root -e "CREATE DATABASE IF NOT EXISTS e_procurement;"
```

If your MySQL user has a password:

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS e_procurement;"
```

### 5. Download Go dependencies

```bash
go mod tidy
```

This will:

- download project dependencies
- generate/update `go.sum`
- prepare the module for build and run

### 6. Create a fresh baseline and seed master data

For the refactored Phase 1 foundation, the safest local setup is a fresh database baseline:

```bash
DB_RESET=true DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go
```

Or:

```bash
make reset-db
```

What this does:

- loads `.env`
- drops the configured database if it already exists
- recreates the database with `utf8mb4`
- creates tables using GORM auto migration
- seeds baseline master data for entity, department, roles, admin user, and `user_roles`
- starts the HTTP server

Seeded defaults:

- entity code: `HO`
- department code: `PROC`
- username: `admin`
- email: `admin@eproc.local`
- password: `Admin123!` unless overridden by `SEED_ADMIN_PASSWORD`
- primary role: `SUPER_ADMIN`

If you only want schema migration without resetting the database, use:

```bash
DB_MIGRATE=true go run ./cmd/api/main.go
```

### 7. Start the API normally

```bash
go run ./cmd/api/main.go
```

If everything is correct, you should see logs similar to:

```text
Database connected successfully
Server starting on :8080 (env: development)
Listening and serving HTTP on :8080
```

### 8. Verify the server health

```bash
curl http://127.0.0.1:8080/health
```

Expected response:

```json
{
  "success": true,
  "message": "OK",
  "data": {
    "service": "e-proc-api",
    "status": "ok",
    "version": "phase-1-foundation"
  },
  "meta": {
    "trace_id": "..."
  }
}
```

## Run the Project

You can run the project in two ways.

### Option 1: Run directly with Go

```bash
go run ./cmd/api/main.go
```

### Option 2: Run using Makefile shortcuts

```bash
make run
```

`make` is optional. It only provides shorter aliases for common commands.

## Database Migration

### Run migration manually

```bash
DB_MIGRATE=true go run ./cmd/api/main.go
```

Or:

```bash
make migrate
```

### Seed baseline master data manually

```bash
DB_SEED=true go run ./cmd/api/main.go
```

Or:

```bash
make seed
```

### Bootstrap schema and master data

```bash
DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go
```

Or:

```bash
make bootstrap
```

### Reset database, migrate, and seed in one command

```bash
DB_RESET=true DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go
```

Or:

```bash
make reset-db
```

### Tables currently created by migration

The current auto migration creates these tables:

- `approval_steps`
- `approval_tasks`
- `audit_logs`
- `bid_items`
- `departments`
- `po_items`
- `pr_items`
- `purchase_orders`
- `purchase_requisitions`
- `rfq_vendors`
- `rfqs`
- `roles`
- `users`
- `vendor_bids`
- `vendors`

Important:

- `DB_MIGRATE=true` creates the schema only.
- `DB_SEED=true` inserts idempotent baseline master data for local development.
- the seeded login uses `admin` with the password from `SEED_ADMIN_PASSWORD`.

## Build and Test

### Run tests

```bash
go test ./...
```

Or:

```bash
make test
```

### Build binary

```bash
go build -o bin/e-proc-api ./cmd/api/main.go
```

Or:

```bash
make build
```

### Tidy dependencies

```bash
go mod tidy
```

Or:

```bash
make tidy
```

## API Base URL

Default local base URL:

```text
http://localhost:8080
```

Primary internal API prefix:

```text
/api/v1/internal
```

## Available Routes

The router currently exposes these routes.

### Public

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/auth/login` | Login |

### Internal protected

These routes require a valid JWT bearer token.

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/api/v1/internal/auth/me` | Get current user info |
| `GET` | `/api/v1/internal/purchase-requests` | List purchase requests |
| `POST` | `/api/v1/internal/purchase-requests` | Create purchase request |
| `GET` | `/api/v1/internal/purchase-requests/:id` | Get purchase request detail |
| `POST` | `/api/v1/internal/purchase-requests/:id/submit` | Submit purchase request |
| `GET` | `/api/v1/internal/rfqs` | List RFQ |
| `POST` | `/api/v1/internal/rfqs` | Create RFQ |
| `GET` | `/api/v1/internal/rfqs/:id` | Get RFQ detail |
| `PATCH` | `/api/v1/internal/rfqs/:id/status` | Update RFQ status |
| `GET` | `/api/v1/internal/purchase-orders` | List purchase orders |
| `POST` | `/api/v1/internal/purchase-orders` | Create purchase order |
| `GET` | `/api/v1/internal/purchase-orders/:id` | Get purchase order detail |
| `PATCH` | `/api/v1/internal/purchase-orders/:id/status` | Update purchase order status |
| `GET` | `/api/v1/internal/vendors` | List vendors |
| `POST` | `/api/v1/internal/vendors` | Create vendor |
| `GET` | `/api/v1/internal/vendors/:id` | Get vendor detail |
| `PUT` | `/api/v1/internal/vendors/:id` | Update vendor |
| `GET` | `/api/v1/internal/approvals/tasks` | Get approval tasks |
| `POST` | `/api/v1/internal/approvals/tasks/:id/approve` | Approve a task |
| `POST` | `/api/v1/internal/approvals/tasks/:id/reject` | Reject a task |

### Additional namespaces scaffolded

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/api/v1/vendor/health` | Vendor namespace health |
| `GET` | `/api/v1/files/health` | Files namespace health |
| `GET` | `/api/v1/reports/health` | Reports namespace health |
| `GET` | `/api/v1/admin/health` | Admin namespace health |
| `GET` | `/api/v1/admin/entities` | List entities |
| `GET` | `/api/v1/admin/entities/:id` | Get entity detail |
| `POST` | `/api/v1/admin/entities` | Create entity (`SUPER_ADMIN`) |
| `GET` | `/api/v1/admin/users` | List users |
| `GET` | `/api/v1/admin/users/:id` | Get user detail |
| `POST` | `/api/v1/admin/users` | Create user |

## Authentication

Authentication uses JWT bearer tokens.

Current authorization baseline:

- `SUPER_ADMIN` has cross-entity access
- `ENTITY_ADMIN` is limited to its own entity scope
- internal procurement routes now enforce role checks
- service-level reads and updates also enforce entity isolation for scoped users

### Login request

```http
POST /api/v1/auth/login
Content-Type: application/json
```

Example payload:

```json
{
  "username": "your-username",
  "password": "your-password"
}
```

Example using `curl`:

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-username",
    "password": "your-password"
  }'
```

Seeded local example:

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'
```

If login succeeds, the API returns:

- `token`
- `expires_at`
- `user`

### Use token for protected routes

```bash
curl http://127.0.0.1:8080/api/v1/internal/auth/me \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Seeded auth baseline

When `DB_SEED=true`, the application bootstraps:

- one internal entity
- one department
- internal and vendor baseline roles
- one internal admin user
- one primary `user_roles` assignment for the admin user

The seeded admin user is intended for local development bootstrap and should be customized with `SEED_ADMIN_PASSWORD`.

## Database Notes

### Verified local MySQL connection example

These values were verified to work in local development:

- host: `localhost`
- port: `3306`
- user: `root`
- password: empty
- database: `e_procurement`

### Direct MySQL connection examples

Connect using TCP:

```bash
mysql -u root -h 127.0.0.1 -P 3306
```

Connect using socket:

```bash
mysql -u root --socket=/tmp/mysql.sock
```

Show databases:

```bash
mysql -u root -e "SHOW DATABASES;"
```

Show application tables:

```bash
mysql -u root -D e_procurement -e "SHOW TABLES;"
```

### Important migration note

The codebase has been refactored toward the BRD, FSD, and TSD foundation. If you already ran an older version of this project before the refactor:

- existing tables may be altered in place during migration
- new TSD-aligned tables such as `entities`, `user_roles`, `budgets`, `purchase_requests`, and `notifications` will be added
- local development is safest with a fresh database if you want a clean schema baseline

Use `DB_RESET=true` or `make reset-db` when you want that clean baseline locally.

## Troubleshooting

### 1. `go: command not found`

Go is not installed or not available in your shell `PATH`.

Check:

```bash
go version
```

Install if needed:

```bash
brew install go
```

### 2. `mysql: command not found`

MySQL client is not installed or not available in your shell `PATH`.

Check:

```bash
mysql --version
```

Install if needed:

```bash
brew install mysql
```

### 3. `Failed to connect to database`

Check:

- MySQL service is running
- database exists
- `.env` credentials are correct
- host and port match your local MySQL instance

Useful commands:

```bash
brew services list
mysql -u root -e "SHOW DATABASES;"
```

### 4. Port `8080` already in use

Find the process:

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
```

Then stop the conflicting process or change `APP_PORT` in `.env`.

### 5. Login fails with `invalid credentials`

Most common causes:

- seed data has not been run yet
- username is wrong
- `SEED_ADMIN_PASSWORD` does not match the password you are trying
- user is inactive

### 6. Tables exist but all endpoints return empty data

This is expected for a fresh database after `DB_MIGRATE=true`.

If you also ran `DB_SEED=true`, only bootstrap master data is inserted. Transactional procurement data remains empty.

## Development Notes

### Main entry point

Application entry point:

- [cmd/api/main.go](/Users/itvico/Dev/e-proc-api/cmd/api/main.go)

### Configuration loader

Environment and config loader:

- [internal/config/config.go](/Users/itvico/Dev/e-proc-api/internal/config/config.go)

### Database connection

Database connection and auto migration:

- [internal/database/database.go](/Users/itvico/Dev/e-proc-api/internal/database/database.go)

### Router

Route definitions:

- [internal/router/router.go](/Users/itvico/Dev/e-proc-api/internal/router/router.go)

### Local setup reference

Additional local setup notes:

- [docs/SETUP_LOCAL.md](/Users/itvico/Dev/e-proc-api/docs/SETUP_LOCAL.md)

### Recommended first checks after setup

```bash
go test ./...
curl http://127.0.0.1:8080/health
mysql -u root -D e_procurement -e "SHOW TABLES;"
```
