package user

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/glue/middleware"
	"github.com/TNAHOM/ATS-system-main/internal/glue/routing"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(
	group *gin.RouterGroup,
	log *zap.Logger,
	userHandler handler.User,
) {
	userRoutes := []routing.Route{
		{
			Method:     http.MethodPost,
			Path:       "/auth/signup",
			Handler:    userHandler.SignUp,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:     http.MethodPost,
			Path:       "/auth/login",
			Handler:    userHandler.LoginUser,
			Middleware: []gin.HandlerFunc{},
		},
		{
			Method:  http.MethodGet,
			Path:    "/user/getAllUsers",
			Handler: userHandler.GetAllUsers,
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
				middleware.AuthUserTypeMiddleware(log, "admin"),
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/user/microservice",
			// Handler: userHandler.GetAllUsers,
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "admin"),
				middleware.ProxyHandler("http://127.0.0.1:8000", log),
			},
		},
	}
	routing.RegisterRoute(group, userRoutes, log)
}
