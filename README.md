# E-Procurement API

Backend API untuk sistem E-Procurement. Repository ini saat ini merepresentasikan Phase 1 foundation: authentication, internal procurement workflow dasar, vendor master data, approval task dasar, admin bootstrap management, dan pondasi schema yang lebih luas untuk fase berikutnya.

## Table of Contents

- [Overview](#overview)
- [Current Status](#current-status)
- [Tech Stack](#tech-stack)
- [Architecture Summary](#architecture-summary)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Local Setup](#local-setup)
- [Run Commands](#run-commands)
- [Bootstrap and Seed Defaults](#bootstrap-and-seed-defaults)
- [API Conventions](#api-conventions)
- [Authentication and Authorization](#authentication-and-authorization)
- [API Routes](#api-routes)
- [Example Payloads](#example-payloads)
- [Database and Migration Notes](#database-and-migration-notes)
- [Troubleshooting](#troubleshooting)
- [Development Notes](#development-notes)
- [References](#references)

## Overview

Repository ini adalah Go monolith dengan pola `cmd/` dan `internal/`.

Secara implementasi, service ini sudah menyediakan:

- login dan JWT-based authentication
- endpoint profil user saat ini
- internal purchase request flow dasar
- internal RFQ flow dasar
- internal purchase order flow dasar
- internal vendor master data
- internal approval task list dan approve/reject action
- admin baseline untuk entity dan user management
- standardized API response dengan `trace_id`
- bootstrap migration dan seed data untuk local development

Repository ini juga sudah punya pondasi model untuk area yang lebih luas seperti:

- budgets
- procurement policy rules
- approval models dan approval matrices
- vendor evaluations dan BAFO
- notifications
- audit logs

Namun, sebagian area tersebut masih berupa schema foundation dan belum seluruhnya punya API atau workflow lengkap.

## Current Status

Status implementasi saat ini mengikuti Phase 1 foundation.

Area yang sudah usable:

- auth login
- auth me
- PR create, list, detail, submit
- RFQ create, list, detail, update status
- PO create, list, detail, update status
- vendor create, list, detail, update
- approval task list, approve, reject
- admin entity list, detail, create
- admin user list, detail, create

Area yang belum lengkap:

- change password
- reset password
- full role assignment management
- department management APIs
- policy and approval master-data APIs
- vendor portal transactional features
- file handling flow
- report/export flow
- queue/worker integration
- complete procurement lifecycle depth

Referensi status detail ada di [IMPLEMENTATION_STATUS.md](/Users/itvico/Dev/e-proc-api/docs/IMPLEMENTATION_STATUS.md).

## Tech Stack

- Go `1.22+`
- Gin untuk HTTP router dan middleware
- GORM untuk ORM
- MySQL sebagai database utama
- JWT via `github.com/golang-jwt/jwt/v5`
- environment loader via `github.com/joho/godotenv`
- bcrypt via `golang.org/x/crypto/bcrypt`

## Architecture Summary

Struktur implementasi mengikuti layering sederhana yang konsisten dengan codebase:

- `cmd/api`
  - application entrypoint
  - load config
  - connect database
  - run migration and seed when requested
  - wire services and handlers
- `internal/config`
  - pembacaan environment variables
  - application config
- `internal/database`
  - database connection
  - database reset/bootstrap
  - auto migration
  - seed baseline master data
- `internal/models`
  - GORM models
  - table mapping
  - domain constants
- `internal/services`
  - business logic dan orchestration
  - data access via GORM
  - entity-scope enforcement untuk beberapa flow
- `internal/handlers`
  - HTTP binding
  - validation handling
  - response mapping
- `internal/middleware`
  - JWT auth
  - request context
  - role guard
  - CORS
- `internal/httpapi`
  - standard success dan error response
- `internal/router`
  - route registration
  - namespace grouping

Desain saat ini belum memisahkan repository layer secara terpisah; service langsung menggunakan `*gorm.DB`. Untuk skala Phase 1 foundation, pendekatan ini masih konsisten dan menjaga diff tetap kecil.

## Project Structure

```text
e-proc-api/
├── cmd/
│   └── api/
│       └── main.go
├── docs/
│   ├── BRD_E-Procurement.md
│   ├── FSD_E-Procurement.md
│   ├── IMPLEMENTATION_STATUS.md
│   └── TSD_E-Procurement.md
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── httpapi/
│   ├── middleware/
│   ├── models/
│   ├── router/
│   └── services/
├── .env.example
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

Sebelum menjalankan project ini, pastikan tersedia:

- Go `1.22` atau lebih baru
- MySQL `8.x` atau versi kompatibel

Perintah verifikasi:

```bash
go version
mysql --version
```

Jika menggunakan Homebrew di macOS:

```bash
brew install go mysql
brew services start mysql
```

## Environment Variables

Mulai dari [.env.example](/Users/itvico/Dev/e-proc-api/.env.example).

### Supported Variables

| Variable | Description | Default |
| --- | --- | --- |
| `APP_PORT` | Port HTTP server | `8080` |
| `APP_ENV` | Environment aplikasi | `development` |
| `DB_HOST` | Host MySQL | `localhost` |
| `DB_PORT` | Port MySQL | `3306` |
| `DB_USER` | User MySQL | `root` |
| `DB_PASSWORD` | Password MySQL | empty |
| `DB_NAME` | Nama database | `e_procurement` |
| `DB_MIGRATE` | Jalankan auto migration saat startup | `false` |
| `DB_RESET` | Drop dan recreate database saat startup | `false` |
| `DB_SEED` | Seed baseline data saat startup | `false` |
| `SEED_ADMIN_PASSWORD` | Password admin hasil seed | `Admin123!` |
| `SEED_ENTITY_CODE` | Entity code seed default | `HO` |
| `SEED_ENTITY_NAME` | Entity name seed default | `Head Office` |
| `SEED_DEPARTMENT_CODE` | Department code seed default | `PROC` |
| `SEED_DEPARTMENT_NAME` | Department name seed default | `Procurement` |
| `JWT_SECRET` | Secret signing JWT | lihat `.env.example` |
| `JWT_EXPIRY_HOURS` | Expiry access token | `24` |
| `JWT_REFRESH_EXPIRY_HOURS` | Expiry refresh token | `168` |

### Example `.env`

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
DB_RESET=false
DB_SEED=false

# Local baseline seed defaults
SEED_ADMIN_PASSWORD=Admin123!
SEED_ENTITY_CODE=HO
SEED_ENTITY_NAME=Head Office
SEED_DEPARTMENT_CODE=PROC
SEED_DEPARTMENT_NAME=Procurement

# JWT
JWT_SECRET=change-this-to-a-long-random-secret
JWT_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_HOURS=168
```

## Local Setup

### 1. Prepare environment file

```bash
cp .env.example .env
```

Sesuaikan credential MySQL bila perlu.

### 2. Ensure MySQL is running

```bash
mysql --version
mysql -u root -e "SELECT VERSION();"
```

Jika database belum ada:

```bash
mysql -u root -e "CREATE DATABASE IF NOT EXISTS e_procurement;"
```

Jika memakai password:

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS e_procurement;"
```

### 3. Download dependencies

```bash
go mod tidy
```

### 4. Bootstrap fresh local database

Opsi paling aman untuk local pertama kali:

```bash
DB_RESET=true DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go
```

Atau:

```bash
make reset-db
```

Flow ini akan:

- memastikan database tersedia
- drop database bila `DB_RESET=true`
- membuat ulang database dengan charset `utf8mb4`
- menjalankan GORM auto migration
- seed baseline entity, department, roles, admin user, dan user-role mapping

### 5. Run the API normally

```bash
go run ./cmd/api/main.go
```

Atau:

```bash
make run
```

### 6. Verify health endpoint

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

## Run Commands

Available targets dari [Makefile](/Users/itvico/Dev/e-proc-api/Makefile):

```bash
make run
make build
make migrate
make seed
make bootstrap
make reset-db
make tidy
make test
```

Equivalent commands:

```bash
go run ./cmd/api/main.go
go build -o bin/e-proc-api ./cmd/api/main.go
go test ./...
go mod tidy
```

## Bootstrap and Seed Defaults

Saat `DB_SEED=true`, aplikasi melakukan seed baseline berikut:

- 1 entity internal
- 1 department pada entity tersebut
- baseline role internal dan vendor
- 1 admin user internal
- 1 primary `user_roles` assignment untuk admin user

Default seeded values:

- entity code: `HO`
- entity name: `Head Office`
- department code: `PROC`
- department name: `Procurement`
- username: `admin`
- email: `admin@eproc.local`
- password: `Admin123!` atau nilai `SEED_ADMIN_PASSWORD`
- primary role: `SUPER_ADMIN`
- scope type: `cross_entity`
- `force_change_password`: `true`

Baseline roles yang di-seed:

- `SUPER_ADMIN`
- `ENTITY_ADMIN`
- `PROCUREMENT_ADMIN`
- `REQUESTER`
- `APPROVER`
- `VENDOR_ADMIN`

Semua seed diupayakan idempotent pada level data master utama melalui pola upsert sederhana.

## API Conventions

### Base URL

```text
http://localhost:8080
```

### Namespace Summary

- `/health`
- `/api/v1/auth/*`
- `/api/v1/internal/*`
- `/api/v1/vendor/*`
- `/api/v1/files/*`
- `/api/v1/reports/*`
- `/api/v1/admin/*`

### Standard Success Response

Semua response sukses mengikuti envelope ini:

```json
{
  "success": true,
  "message": "OK",
  "data": {},
  "meta": {
    "trace_id": "..."
  }
}
```

Untuk create operation, message menjadi `"Created"`.

### Standard Error Response

Response error mengikuti bentuk ini:

```json
{
  "success": false,
  "message": "Validation failed",
  "error_code": "VALIDATION_ERROR",
  "errors": [
    {
      "field": "body",
      "message": "..."
    }
  ],
  "meta": {
    "trace_id": "..."
  }
}
```

### Trace ID

Setiap request mendapat `trace_id` melalui middleware request context. Nilai ini muncul pada response sukses maupun error dan berguna untuk tracing log/debugging.

### Pagination

Beberapa endpoint list menggunakan query parameter:

- `page`
- `page_size`

Default:

- `page=1`
- `page_size=20`

Batas maksimum `page_size` saat ini adalah `100`.

## Authentication and Authorization

### Login Endpoint

```http
POST /api/v1/auth/login
Content-Type: application/json
```

Request body:

```json
{
  "username": "admin",
  "password": "Admin123!"
}
```

Example:

```bash
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'
```

Successful login returns:

- `token`
- `expires_at`
- `user`

`user` saat ini dapat memuat:

- `id`
- `entity_id`
- `username`
- `email`
- `full_name`
- `role_code`
- `role_name`
- `scope_type`
- `department_name`

### Authenticated User Endpoint

```bash
curl http://127.0.0.1:8080/api/v1/internal/auth/me \
  -H "Authorization: Bearer <token>"
```

Response `data` memuat:

- `user_id`
- `entity_id`
- `username`
- `role_code`
- `role_name`
- `scope_type`
- `subject_type`

### Authorization Baseline

Current authorization baseline di code:

- `SUPER_ADMIN`
  - akses lintas entity
  - bisa create entity
  - bisa manage admin namespace secara penuh dalam scope saat ini
- `ENTITY_ADMIN`
  - terbatas ke entity sendiri
  - bisa akses admin namespace dalam entity scope
- `PROCUREMENT_ADMIN`
  - bisa mengakses namespace procurement internal yang relevan
- `REQUESTER`
  - bisa create dan melihat purchase request dalam namespace internal
- `APPROVER`
  - bisa mengakses approval task internal

Entity scope enforcement diterapkan di:

- detail read tertentu
- update status tertentu
- submit/approval flows tertentu
- user creation rules

## API Routes

### Public Routes

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/auth/login` | Login |

### Internal Protected Routes

Perlu bearer token dan role yang sesuai.

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/api/v1/internal/auth/me` | Current authenticated user |
| `GET` | `/api/v1/internal/purchase-requests` | List purchase requests |
| `POST` | `/api/v1/internal/purchase-requests` | Create purchase request |
| `GET` | `/api/v1/internal/purchase-requests/:id` | Purchase request detail |
| `POST` | `/api/v1/internal/purchase-requests/:id/submit` | Submit purchase request |
| `GET` | `/api/v1/internal/rfqs` | List RFQs |
| `POST` | `/api/v1/internal/rfqs` | Create RFQ |
| `GET` | `/api/v1/internal/rfqs/:id` | RFQ detail |
| `PATCH` | `/api/v1/internal/rfqs/:id/status` | Update RFQ status |
| `GET` | `/api/v1/internal/purchase-orders` | List purchase orders |
| `POST` | `/api/v1/internal/purchase-orders` | Create purchase order |
| `GET` | `/api/v1/internal/purchase-orders/:id` | Purchase order detail |
| `PATCH` | `/api/v1/internal/purchase-orders/:id/status` | Update purchase order status |
| `GET` | `/api/v1/internal/vendors` | List vendors |
| `POST` | `/api/v1/internal/vendors` | Create vendor |
| `GET` | `/api/v1/internal/vendors/:id` | Vendor detail |
| `PUT` | `/api/v1/internal/vendors/:id` | Update vendor |
| `GET` | `/api/v1/internal/approvals/tasks` | My approval tasks |
| `POST` | `/api/v1/internal/approvals/tasks/:id/approve` | Approve task |
| `POST` | `/api/v1/internal/approvals/tasks/:id/reject` | Reject task |

### Additional Protected Namespace Routes

| Method | Route | Description |
| --- | --- | --- |
| `GET` | `/api/v1/vendor/health` | Vendor namespace health |
| `GET` | `/api/v1/files/health` | Files namespace health |
| `GET` | `/api/v1/reports/health` | Reports namespace health |
| `GET` | `/api/v1/admin/health` | Admin namespace health |
| `GET` | `/api/v1/admin/entities` | List entities |
| `GET` | `/api/v1/admin/entities/:id` | Entity detail |
| `POST` | `/api/v1/admin/entities` | Create entity, `SUPER_ADMIN` only |
| `GET` | `/api/v1/admin/users` | List users |
| `POST` | `/api/v1/admin/users` | Create user |
| `GET` | `/api/v1/admin/users/:id` | User detail |

### Common Query Parameters

List endpoint yang saat ini mendukung filter/query:

- purchase requests
  - `page`
  - `page_size`
  - `status`
  - `department_code`
- RFQs
  - `page`
  - `page_size`
  - `status`
- purchase orders
  - `page`
  - `page_size`
  - `status`
  - `vendor_id`
- vendors
  - `page`
  - `page_size`
  - `active_only`
- users
  - `entity_id`
  - `status`

Beberapa filter entity diterapkan otomatis berdasarkan claim/token caller.

## Example Payloads

Contoh di bawah ini mengikuti request shape yang saat ini digunakan oleh service/handler.

### Create Purchase Request

```http
POST /api/v1/internal/purchase-requests
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "title": "Pengadaan Laptop Operasional",
  "description": "Pengadaan laptop untuk tim procurement",
  "department_code": "PROC",
  "procurement_type": "goods",
  "routine_type": "non_routine",
  "budget_status": "within_budget",
  "need_date": "2026-04-10T00:00:00Z",
  "items": [
    {
      "item_name": "Laptop Business",
      "specification": "RAM 16GB SSD 512GB",
      "qty": 5,
      "uom": "unit",
      "estimated_unit_price": 15000000
    }
  ]
}
```

### Submit Purchase Request

```http
POST /api/v1/internal/purchase-requests/:id/submit
Authorization: Bearer <token>
```

Behavior penting:

- hanya PR dengan status `Draft` atau `Revised` yang dapat di-submit
- submit akan mengubah status menjadi `Pending Approval`
- submit juga membuat `approval_tasks` record

### Create RFQ

```http
POST /api/v1/internal/rfqs
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "pr_id": 1,
  "title": "RFQ Laptop Operasional",
  "technical_requirement": "Minimal RAM 16GB",
  "commercial_requirement": "Garansi resmi 3 tahun",
  "minimum_vendor_count": 2,
  "deadline_at": "2026-04-15T00:00:00Z",
  "vendor_ids": [1, 2]
}
```

Behavior penting:

- PR harus berada dalam entity scope caller
- `minimum_vendor_count` otomatis menjadi `1` bila dikirim `0` atau negatif

### Update RFQ Status

```http
PATCH /api/v1/internal/rfqs/:id/status
Authorization: Bearer <token>
Content-Type: application/json
```

Status values saat ini tidak dibatasi ketat di handler, tetapi model menyediakan status constants seperti:

- `Created`
- `Published`
- `Vendor Submission`
- `Closed`
- `Reopened`
- `Evaluation`
- `BAFO`
- `Vendor Selected`
- `Cancelled`

### Create Purchase Order

```http
POST /api/v1/internal/purchase-orders
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "pr_id": 1,
  "rfq_id": 1,
  "vendor_id": 1,
  "po_date": "2026-04-20T00:00:00Z",
  "expected_date": "2026-04-30T00:00:00Z",
  "delivery_address": "Jakarta Head Office",
  "payment_terms": "30 days after invoice",
  "notes": "Handle with care",
  "items": [
    {
      "pr_item_id": 1,
      "item_name": "Laptop Business",
      "specification": "RAM 16GB SSD 512GB",
      "qty": 5,
      "uom": "unit",
      "unit_price": 14500000
    }
  ]
}
```

Behavior penting:

- `vendor_id`, `po_date`, `delivery_address`, dan minimal satu item wajib ada
- bila `pr_id` atau `rfq_id` dikirim, record tersebut harus masih berada dalam entity scope caller
- currency default saat ini adalah `IDR`
- status awal saat create adalah `Draft`

### Create Vendor

```http
POST /api/v1/internal/vendors
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "vendor_name": "PT Vendor Teknologi",
  "tax_id": "01.234.567.8-999.000",
  "email": "vendor@example.com",
  "phone": "0211234567",
  "address": "Jakarta"
}
```

Behavior penting:

- code vendor digenerate otomatis dengan format `V-0001`
- status default vendor baru:
  - `approved_status=approved`
  - `blacklist_status=false`
  - `eligibility_status=eligible`

### Create Entity

```http
POST /api/v1/admin/entities
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "entity_code": "SUB1",
  "entity_name": "Subsidiary 1",
  "entity_type": "subsidiary",
  "governance_mode": "entity_only",
  "status": "active"
}
```

Defaulting behavior:

- `entity_type` default `subsidiary`
- `governance_mode` default `entity_only`
- `status` default `active`

### Create User

```http
POST /api/v1/admin/users
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "entity_id": 1,
  "department_id": 1,
  "full_name": "Procurement Staff",
  "email": "proc.staff@example.com",
  "username": "procstaff",
  "password": "StrongPass123",
  "role_code": "REQUESTER",
  "scope_type": "own_entity",
  "status": "active",
  "force_change_password": true
}
```

Behavior penting:

- password minimal 8 karakter
- role harus ada dan aktif
- `ENTITY_ADMIN` tidak boleh membuat user di luar entity miliknya
- `ENTITY_ADMIN` juga tidak boleh assign scope `cross_entity`

## Database and Migration Notes

### Auto-Migrated Main Tables

Migration saat ini membuat atau menyesuaikan area tabel berikut:

- `entities`
- `departments`
- `roles`
- `users`
- `user_roles`
- `delegate_approvers`
- `vendors`
- `vendor_users`
- `vendor_blacklists`
- `reference_prices`
- `budgets`
- `procurement_policy_rules`
- `approval_models`
- `approval_matrices`
- `purchase_requests`
- `purchase_request_items`
- `pr_attachments`
- `pr_approvals`
- `rfqs`
- `rfq_vendors`
- `quotations`
- `quotation_items`
- `vendor_evaluations`
- `bafo_rounds`
- `vendor_selections`
- `direct_appointments`
- `purchase_orders`
- `purchase_order_items`
- `po_approvals`
- `vendor_confirmations`
- `approval_tasks`
- `notifications`
- `audit_logs`
- `app_logs`

### Migration Behavior

Catatan penting:

- migration menggunakan GORM `AutoMigrate`
- foreign key constraint creation saat migration dinonaktifkan
- strategi ini dipakai untuk mengurangi konflik dengan schema lokal lama yang pernah ada

### Existing Database Caution

Jika sebelumnya pernah menjalankan versi lama project ini:

- beberapa tabel lama bisa masih tersisa
- beberapa tabel akan berubah in place
- hasil schema bisa campuran antara baseline lama dan foundation baru

Untuk local development, baseline paling bersih adalah:

```bash
DB_RESET=true DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go
```

### Quick Database Checks

```bash
mysql -u root -e "SHOW DATABASES;"
mysql -u root -D e_procurement -e "SHOW TABLES;"
mysql -u root -D e_procurement -e "SELECT username, email, status FROM users;"
mysql -u root -D e_procurement -e "SELECT role_code, role_name FROM roles;"
```

## Troubleshooting

### `go: command not found`

Go belum ter-install atau belum masuk ke `PATH`.

Check:

```bash
go version
```

### `mysql: command not found`

MySQL client belum tersedia di shell.

Check:

```bash
mysql --version
```

### Failed to connect to database

Periksa:

- service MySQL berjalan
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, dan `DB_NAME` sudah benar
- database user memiliki akses ke database target

Useful commands:

```bash
mysql -u root -e "SHOW DATABASES;"
mysql -u root -e "SELECT VERSION();"
```

### Port `8080` already in use

Check process:

```bash
lsof -nP -iTCP:8080 -sTCP:LISTEN
```

Lalu ganti `APP_PORT` atau stop process yang bentrok.

### Login fails with `invalid credentials`

Periksa:

- apakah `DB_SEED=true` pernah dijalankan
- username/password sesuai seed atau user yang dibuat
- user status masih `active`

### Login fails because account is locked

Service auth saat ini mengenali kondisi `locked_until` pada user. Jika field tersebut terisi dengan waktu yang masih aktif, login akan ditolak.

### Endpoint returns validation error

Penyebab umum:

- JSON body tidak valid
- field required belum dikirim
- tipe data tidak sesuai
- email tidak valid
- password kurang dari 8 karakter pada create user
- list item kosong pada create PR atau create PO

### Data list kosong setelah migration

Itu normal jika hanya menjalankan:

```bash
DB_MIGRATE=true go run ./cmd/api/main.go
```

Migration hanya membuat schema. Untuk baseline master data, jalankan juga seed:

```bash
DB_SEED=true go run ./cmd/api/main.go
```

Atau gunakan bootstrap/reset-db.

## Development Notes

Entry points penting:

- [main.go](/Users/itvico/Dev/e-proc-api/cmd/api/main.go)
- [config.go](/Users/itvico/Dev/e-proc-api/internal/config/config.go)
- [database.go](/Users/itvico/Dev/e-proc-api/internal/database/database.go)
- [bootstrap.go](/Users/itvico/Dev/e-proc-api/internal/database/bootstrap.go)
- [router.go](/Users/itvico/Dev/e-proc-api/internal/router/router.go)
- [response.go](/Users/itvico/Dev/e-proc-api/internal/httpapi/response.go)
- [auth.go](/Users/itvico/Dev/e-proc-api/internal/services/auth.go)

Recommended local verification:

```bash
go test ./...
curl http://127.0.0.1:8080/health
curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"Admin123!"}'
mysql -u root -D e_procurement -e "SHOW TABLES;"
```

## References

- [BRD_E-Procurement.md](/Users/itvico/Dev/e-proc-api/docs/BRD_E-Procurement.md)
- [FSD_E-Procurement.md](/Users/itvico/Dev/e-proc-api/docs/FSD_E-Procurement.md)
- [TSD_E-Procurement.md](/Users/itvico/Dev/e-proc-api/docs/TSD_E-Procurement.md)
- [IMPLEMENTATION_STATUS.md](/Users/itvico/Dev/e-proc-api/docs/IMPLEMENTATION_STATUS.md)
