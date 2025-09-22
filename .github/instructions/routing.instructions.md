# Routing Instructions

## How to edit this file
- Keep routes grouped by resource
- Use RegisterRoute consistently
- Attach AuthMiddleware when needed

---

## Rules for routing
- Define all routes inside `internal/glue/routing/`.
- Use the `Route` struct with `Method`, `Path`, `Handler`, `Middleware`.
- Group routes by resource (e.g., user, auth).
- Apply `AuthMiddleware(log)` on protected routes.
- For proxy routes, use `middleware.ProxyHandler(target, log)`.

## Example
```go
func Init(
    group *gin.RouterGroup,
    log *zap.Logger,
    userHandler handler.User,
) {
    userRoutes := []routing.Route{
        {
            Method:  http.MethodPost,
            Path:    "/auth/signup",
            Handler: userHandler.SignUp,
        },
        {
            Method:  http.MethodGet,
            Path:    "/user/getAllUsers",
            Handler: userHandler.GetAllUsers,
            Middleware: []gin.HandlerFunc{
                middleware.AuthMiddleware(log),
            },
        },
    }
    routing.RegisterRoute(group, userRoutes, log)
}
