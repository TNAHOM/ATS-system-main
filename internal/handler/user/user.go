package user

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct {
	log        *zap.Logger
	userModule module.User
}

func Init(log *zap.Logger, userModule module.User) handler.User {
	return &user{log: log, userModule: userModule}
}

// SignUp godoc
// @Summary Sign up user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body dto.CreateUserRequest true "Sign up payload"
// @Success 200 {object} dto.EnvelopeCreateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/signup [post]
func (u *user) SignUp(ctx *gin.Context) {
	var userModel dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userModel); err != nil {
		u.log.Error(err.Error(), zap.Any("request", userModel))
		ctx.JSON(http.StatusBadRequest, dto.Envelope[any]{Error: err.Error()})
		return
	}

	res, err := u.userModule.CreateUser(ctx, userModel)
	if err != nil {
		u.log.Error(err.Error(), zap.Any("request", userModel))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, dto.Envelope[dto.CreateUserResponse]{Data: res})
}

// LoginUser godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body dto.LoginUserRequest true "Login payload"
// @Success 200 {object} dto.EnvelopeLoginUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (u *user) LoginUser(ctx *gin.Context) {
	var loginModel dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&loginModel); err != nil {
		u.log.Error(err.Error(), zap.Any("request", loginModel))
		ctx.JSON(http.StatusBadRequest, dto.Envelope[any]{Error: err.Error()})
		return
	}

	loginRes, err := u.userModule.LoginUser(ctx, loginModel)
	if err != nil {
		u.log.Error(err.Error(), zap.Any("request", loginModel))
		ctx.JSON(http.StatusInternalServerError, dto.Envelope[any]{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Envelope[dto.LoginUserResponse]{Data: loginRes})
}

// GetAllUsers godoc
// @Summary List all users (admin only)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.EnvelopeGetAllUsersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user/getAllUsers [get]
func (u *user) GetAllUsers(ctx *gin.Context) {
	users, err := u.userModule.GetAllUsers(ctx)
	if err != nil {
		u.log.Error("failed to get users", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, dto.Envelope[any]{Error: "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, dto.Envelope[[]dto.GetAllUsers]{Data: users})
}
