package aimicroservice

import (
	"net/http"
	"os"

	"github.com/TNAHOM/ATS-system-main/internal/glue/middleware"
	"github.com/TNAHOM/ATS-system-main/internal/glue/routing"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(group *gin.RouterGroup, log *zap.Logger) {
	aiBaseUrl := os.Getenv("AI_SERVICE_URL")
	aiMicroserviceRoutes := []routing.Route{
		{
			Method:     http.MethodPost,
			Path:       "/resumes/upload",
			Handler:    middleware.ProxyHandler(aiBaseUrl+"/resumes/upload", log),
			Middleware: []gin.HandlerFunc{
				// middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "admin"),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/resumes/application/:applicationID",
			Handler: middleware.ProxyHandler(aiBaseUrl+"/resumes/application/:applicationID", log),
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "admin"),
			},
		},
	}
	routing.RegisterRoute(group, aiMicroserviceRoutes, log)
}