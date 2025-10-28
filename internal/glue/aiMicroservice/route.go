package aimicroservice

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/glue/middleware"
	"github.com/TNAHOM/ATS-system-main/internal/glue/routing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(group *gin.RouterGroup, log *zap.Logger) {
	aiMicroserviceRoutes := []routing.Route{
		{
			Method:     http.MethodPost,
			Path:       "/resumes/upload",
			Handler:    middleware.ProxyHandler("http://127.0.0.1:8000/resumes/upload", log),
			Middleware: []gin.HandlerFunc{
				// middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "admin"),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/resumes/application/:applicationID",
			Handler: middleware.ProxyHandler("http://127.0.0.1:8000/resumes/application/:applicationID", log),
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "admin"),
			},
		},
	}
	routing.RegisterRoute(group, aiMicroserviceRoutes, log)
}
