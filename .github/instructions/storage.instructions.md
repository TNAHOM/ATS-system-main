# Storage Instructions

## How to edit this file
- Storage = GORM only
- Must map DTO ↔ Model ↔ DTO
- Always use db.WithContext

---

## Rules for storage
- All DB operations must use `db.WithContext(ctx)`.
- Storage layer returns DTOs, not models.
- Never leak raw DB errors; wrap them.
- Inject `*zap.Logger` for error logging.

## Example
```go
type UserRepository struct {
    db  *gorm.DB
    log *zap.Logger
}

func (r *UserRepository) Create(ctx context.Context, user model.User) error {
    if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
        r.log.Error("db create failed", zap.Error(err))
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}
