
# Module Instructions

## How to edit this file
- Keep modules pure: no Gin, no swagger, no HTTP
- Modules = business logic
- Always work with DTOs

---

## Rules for modules
- Modules implement **business logic only**.
- Accept DTOs, return DTOs.
- Never import `gin` or framework packages.
- Always inject `*zap.Logger` into module structs.
- Call storage interfaces, never DB directly.

## Example
```go
type UserModule struct {
    repo storage.UserRepository
    log  *zap.Logger
}

func (m *UserModule) SignUp(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
    user := model.User{
        Username: req.Username,
        Email:    req.Email,
    }
    if err := m.repo.Create(ctx, user); err != nil {
        m.log.Error("failed to create user", zap.Error(err))
        return dto.CreateUserResponse{}, fmt.Errorf("create user: %w", err)
    }
    return dto.CreateUserResponse{ID: user.ID}, nil
}
