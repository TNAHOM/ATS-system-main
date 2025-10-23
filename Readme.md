# ATS-system-main

A modular Applicant Tracking System (ATS) backend written in Go, showcasing a clean DDD + hexagonal architecture with Gin, GORM (PostgreSQL), JWT auth, zap logging, and pgvector-powered AI embeddings for job posts.

## Highlights

- Domain-first design: Handlers → Modules → Storage separation (no handler talks to DB directly).
- Consistent DTOs for requests/responses and a generic Envelope[T] response wrapper.
- JWT auth (HMAC-SHA256) with access and refresh tokens; bcrypt password hashing.
- Auth middleware that validates Bearer tokens and places claims in request context.
- Job post embeddings via Google GenAI and pgvector (3072-dim vectors) for future semantic search/matching.
- Zap logging across layers; errors logged with context.

## Tech stack

- Language: Go (module requires go 1.24 as per `go.mod`)
- Web: Gin (`github.com/gin-gonic/gin`)
- Auth: `github.com/golang-jwt/jwt`, bcrypt (x/crypto)
- DB: PostgreSQL + GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`)
- Vectors: `pgvector` extension and `pgvector-go`
- Logging: `go.uber.org/zap`
- Env: `github.com/joho/godotenv`

## Architecture overview

The project enforces a hexagonal layout:

- Handlers (HTTP): validate/bind, call modules, return `Envelope[T]` JSON
- Modules (domain): business logic only; consume storage interfaces and platform services
- Storage (infra): GORM implementations mapping DTO ↔ Model ↔ DTO using `db.WithContext(ctx)`
- Platform: cross-cutting (encryption/JWT, AI embeddings, response helpers)

Key files:

- Entry: `cmd/main.go` → `initiator/initiator.go`
- Wiring: `initiator/{log,config,db,platform,persistance,module,handler,route,syncDB}.go`
- Routing: `internal/glue/{routing,routes per domain}` + `internal/glue/middleware/auth.go`
- Handlers: `internal/handler/{user,jobPost}`
- Modules: `internal/module/{user,jobPost}`
- Storage: `internal/storage/{user,jobPost}`
- DTOs: `internal/constants/dto` and Models: `internal/constants/model`

## Requirements

- Go 1.24+ (per `go.mod`)
- PostgreSQL 14+ with `pgvector` extension installed
- A Google API key for embeddings (if creating/updating job posts)

Enable pgvector in your database (once):

```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

The schema is created at startup via GORM AutoMigrate. Job post tables declare three `vector(3072)` columns.

## Environment variables

These are loaded from `.env` by `initiator/config.go`:

- `HOST` – HTTP bind host (e.g., 0.0.0.0)
- `PORT` – HTTP port (e.g., 8080)
- `DB_URL` – Postgres DSN (e.g., `postgres://user:pass@localhost:5432/ats?sslmode=disable`)
- `SECRET_KEY` – HMAC secret for JWT signing/validation
- `GOOGLE_API_KEY` – required for embeddings (job post create/update)

Example `.env`:

```env
HOST=0.0.0.0
PORT=8080
DB_URL=postgres://postgres:postgres@localhost:5432/ats?sslmode=disable
SECRET_KEY=change-me
GOOGLE_API_KEY=your-google-api-key
```

## Run locally

1) Install dependencies (Go modules are automatic on build).

2) Ensure PostgreSQL is running and `pgvector` is installed (see SQL above).

3) Create `.env` with the variables above.

4) Start the server:

```bash
go run ./cmd
```

The API will listen on `http://HOST:PORT` and mount all routes under `/api` (see `initiator/initiator.go`).

## API endpoints

All responses follow the generic envelope:

```json
{
      "data": { /* type varies */ },
      "error": "" /* omitted on success */
}
```

### Auth and Users

- POST `/api/auth/signup`
      - Body (CreateUserRequest):
            - first_name, last_name, email, phone, password, user_type
      - On success (CreateUserResponse): returns user info plus `token` and `refresh_token`.

- POST `/api/auth/login`
      - Body: { email, password }
      - On success (LoginUserResponse): returns user info plus fresh `token` and `refresh_token`.

- GET `/api/user/getAllUsers` (protected; requires Bearer token AND `user_type=admin`)
      - Returns: array of users (id, first_name, last_name, email, phone, user_type)

Auth details:

- Send `Authorization: Bearer <access_token>`
- Claims are injected into Gin context and `context.Context` for downstream modules.

### Job posts

All job post endpoints require a valid Bearer token.

- POST `/api/jobPost/create`
      - Body (CreateJobPostRequest): title, description, deadline, responsibilities[], requirements[]
      - The service generates 3072-dim embeddings via Google GenAI and persists them with pgvector.
      - Response (CreateJobPostResponse): persisted job post core fields.

- GET `/api/jobPost/getAllJobPosts`
      - Response: array of job posts (id, title, description, responsibilities, requirements, user_id, deadline)

- PATCH `/api/jobPost/update/:id`
      - Partial update; only provided fields are changed. If description/responsibilities/requirements change, embeddings are regenerated.
      - Response: full updated record (UpdateJobPostResponse)

- DELETE `/api/jobPost/:id`
      - Soft-deletes the record. Returns `{ "data": { "deleted": true } }` on success.

### Example (curl)

```bash
# Sign up
curl -sS -X POST http://localhost:8080/api/auth/signup \
      -H "Content-Type: application/json" \
      -d '{
            "first_name":"Jane",
            "last_name":"Doe",
            "email":"jane@example.com",
            "phone":"+10000000000",
            "password":"P@ssw0rd",
            "user_type":"admin"
      }'

# Login
TOKEN=$(curl -sS -X POST http://localhost:8080/api/auth/login \
      -H "Content-Type: application/json" \
      -d '{"email":"jane@example.com","password":"P@ssw0rd"}' | jq -r .data.token)

# Admin-only
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/user/getAllUsers
```

Notes:

- The login response contains `data.token` and `data.refresh_token`.
- The admin-only endpoint uses both `AuthMiddleware` and `AuthUserTypeMiddleware("admin")`.

### Swagger/OpenAPI

- swag init -g cmd/main.go -o docs

## Error handling and logging

- Handlers use `dto.Envelope[T]` for success and error responses.
- Some handlers use `platform/response.SendError` to format validation errors (via go-playground/validator).
- Zap logger (`*zap.Logger`) is injected and used in all layers; errors are logged with context.

## Development notes

- AutoMigrate runs at startup for `User` and `JobPost` models (`initiator/syncDB.go`).
- Storage implementations always call `db.WithContext(ctx)` and map models ↔ DTOs.
- Do not return raw DB errors to handlers; log and propagate wrapped errors where appropriate.

## Contributing

1) Open an issue describing the change.
2) Follow the DDD conventions and DTO/Envelope patterns.
3) Keep handlers thin and push logic into modules; never call DB from handlers.

---

## Project TODO

- [ ] Add Swagger minimal annotations and generate OpenAPI docs for all routes.
- [ ] Provide Docker and docker-compose (app + PostgreSQL with pgvector) for one-command setup.
- [ ] Implement token refresh endpoint and refresh token rotation/revocation.
- [ ] Add pagination and filtering for list endpoints (users, job posts).
- [ ] Introduce structured error types → HTTP code mapping in one place.
- [ ] Add CI (lint, vet, tests) and pre-commit hooks (gofmt, goimports).
- [ ] Harden validation for user and job post DTOs; return consistent validation messages.
- [ ] Add health/ready endpoints and basic metrics.
- [ ] Optional: role-based access middleware improvements and per-route permissions.
