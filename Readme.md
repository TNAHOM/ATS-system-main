# Simple JWT Authentication (example)

This small project demonstrates a minimal JWT-based authentication flow in Go and a tiny middleware setup that protects routes by validating tokens and enforcing user type.

## What this repository shows

- Token generation using HMAC-SHA256 JWTs with a 24-hour access token and a 7-day refresh token.
- Password hashing with bcrypt.
- Token validation middleware that extracts Bearer tokens and puts `user_id` and `user_type` into the Gin context.
- A protected route example that requires the `admin` user type.

## Quick overview

- Encryption implementation: `platform/encryption/encryption.go` — GenerateToken, ValidateToken, HashPassword, VerifyPassword.
- Middleware: `internal/glue/middleware/auth.go` — checks Authorization header, validates token, sets `user_id` and `user_type` in context.
- Routes: `internal/glue/user/user.go` — public endpoints:

  - POST `/api/auth/signup`
  - POST `/api/auth/login`
  - GET `/api/user/getAllUsers` (protected, requires `Authorization: Bearer <token>` and `user_type = admin`)

## Environment variables

- `SECRET_KEY` — HMAC secret used to sign and validate JWTs (required).
- `HOST`, `PORT` — server bind address (used in `initiator/initiator.go`).

## How tokens work (short)

- Access token: expires in 24 hours and carries user claims (id, email, name, user type).
- Refresh token: expires in 7 days.
- Tokens are signed with `SECRET_KEY`.

## Example (curl)

- Sign up / login: send JSON to `/api/auth/login` or `/api/auth/signup`.
- Access protected route:

```bash
# replace TOKEN with the access token returned from login
curl -H "Authorization: Bearer TOKEN" http://localhost:8080/api/user/getAllUsers
```

## Notes and next steps

- This example uses environment-based secret management. For production, use a secure secret store and rotate keys.
- Consider adding token revocation (store active refresh tokens), HTTPS, rate limiting, and stricter error handling.

That's it — a compact demonstration of JWT auth and minimal middleware in a Gin-based Go app.

## Recommended folder structure

```text
├── bin                                 
├── cmd                                 # Packages that provide support for a specific program that is being built, this contains cmd for REST, gRPC, etc
├── initiator                           # initialize any domain here that will be called inside cmd, so, no spaghetti init on cmd main package
├── internal                            # Contains all application packages (all functions must have return)
      ├── constants                     # For declaring all constants variable based on it's entity, declaring with full name descripting it
            ├── model                   # model is where the declaration of structs are being written
            ├── query                   # where query is being declared for storage sql, and it's better using sqlx naming system for passing the parameters
            ├── state                   # contains all constants variable (string, int, etc..) like status, or constant strings
      ├── glue                          # is just another name of middleware
            ├── routing                 # where the routing for handlers are assigned based on method and url
      ├── handler                       # any protocol interface that needs handler is being declare here, whether for cleaning input data, error payload handling, etc..
            ├── rest                    # handler for rest API technology
      ├── module                        # Domain packages are here, contains business logic and interfaces that belong to each domain, and no usecase calling other usecase
            ├── user                    # only sample domain, user package which handler for user business logic
                ├── initiator.go        # this is the only file to declare interface methods from storage and repository. also where to put func init the package.
                ....                    # name the file based on the usecase that may contains multiple acticity; example: login usecase for only one activity, profile contains edit profile and show profile
      ├── repository                    # this may a bit weird for you, this package uses for data storing logic, it is an optional if your domain only saves data to one db, but it's different when a domain uses multiple storage, for example caching and multiple persistences
      ├── storage                       # this is where you put the data storing code. whether persistence like postgresql, monggodb, etc. and caching like redis, etc. 
            ├── cache                   # package where caching storing code is written based on its domain.
            ├── persistence             # like its name, this package contains storing code for SQL or noSQL db
├── mocks                               # this contains all mock from any source. it contains some rules. it must be a generate files and the package name must be mock_(source)
      ├── psql                          # this is mock for psql contains package mock_psql, even tho it is written manually, it must have // Code generated manually. DO NOT EDIT.
      ├── redis                         # this is mock for redis contains package mock_redis, even tho it is written manually, it must have // Code generated manually. DO NOT EDIT.
      ├── user                          # this is generated by golang mock official generator
├── platform                            # external app for all uses (in here you can make function without return)
      ├── encryption                    # you can have any encryption method here
      ├── postgres                      # contains functions to open database postgres connections, with mutiple servers can be added, like db master and/or slave
      ├── redis                         # contains functions to open database redis connection, currently it uses only one connection, but this can be adjust just like the postgres connections
      ├── routers                       # contains functions to serve the HTTP Listener using all registered URLs with the handlers
├── static                              # to store any static file like pic, html, etc 
```
