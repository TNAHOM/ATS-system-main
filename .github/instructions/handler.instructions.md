# Handler Instructions

## How to edit this file
- Keep handler examples aligned with DTO + Envelope + zap logger
- Always use swagger annotations minimally
- Handlers should remain thin: delegate to modules

---

## Rules for handlers
- Always accept and validate input using `ctx.ShouldBindJSON(&dto)`.
- Always return responses using `Envelope[T]`.
- Inject `*zap.Logger` in every handler struct.
- Handlers must call **modules only**, never storage directly.
- Add minimal Swagger annotations for every handler.

## Example
```go
// @Summary Sign up user
// @Success 200 {object} dto.Envelope[dto.CreateUserResponse]
// @Failure 400 {object} dto.Envelope[any]
// @Router /auth/signup [post]
func (h *userHandler) SignUp(ctx *gin.Context) {
    var req dto.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        h.log.Error("invalid request", zap.Error(err))
        ctx.JSON(http.StatusBadRequest, dto.Envelope[any]{Error: "invalid input"})
        return
    }

    res, err := h.module.SignUp(ctx, req)
    if err != nil {
        h.log.Error("signup failed", zap.Error(err))
        ctx.JSON(http.StatusInternalServerError, dto.Envelope[any]{Error: "internal error"})
        return
    }

    ctx.JSON(http.StatusOK, dto.Envelope[dto.CreateUserResponse]{Data: res})
}
