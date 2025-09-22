# Copilot Instructions (Repo-wide)

## How to edit this file
- Update examples when new patterns are introduced
- Keep envelope, DTO, and zap logger rules consistent
- Do not remove core conventions (DDD, hexagonal, Gin + GORM)

---

## General Go conventions
- Always use idiomatic Go style (`gofmt`, `goimports`).
- Use `if err != nil { ... }` style error handling.
- Wrap errors with `%w` and prefer `errors.Is/As` checks.
- No hardcoded config values: load via `env` or config layer.

## Project architecture
- Enforce strict **DDD + hexagonal structure**:
  - Handlers → Modules → Storage.
  - Handlers never call storage directly.
- Use DTOs for all input/output (never expose raw models).
- Models always live in `internal/constants/model`.

## Logging
- Every handler/module/storage must inject and use `*zap.Logger`.
- Log errors with context (`log.Error("msg", zap.Error(err))`).
- Do not ignore errors silently.

## Error handling
- All HTTP responses must use the generic `Envelope[T]` type.
- Distinguish error types with proper HTTP codes:
  - 400 for validation errors
  - 401 for unauthorized
  - 403 for forbidden
  - 404 for not found
  - 500 for internal errors

## Swagger
- Use **minimal annotations**: `@Summary`, `@Success`, `@Failure`, `@Router`.
- Do not generate long descriptions unless explicitly asked.

## Database
- Always use `db.WithContext(ctx)`.
- Storage layer maps DTO ↔ Model ↔ DTO.
- Never return raw DB errors directly; wrap and propagate.

## Middleware & Auth
- If a route requires protection, always apply `AuthMiddleware(log)`.
- JWT validation lives in `platform/encryption`.
