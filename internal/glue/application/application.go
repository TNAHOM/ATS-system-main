package application

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/glue/routing"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(group *gin.RouterGroup, log *zap.Logger, applicationHandler handler.Application) {
	applicationRoutes := []routing.Route{
		{
			Method:     http.MethodGet,
			Path:       "/applications/jobPostByID/:jobPostID",
			Handler:    applicationHandler.GetApplicationsByJobPostID,
			Middleware: []gin.HandlerFunc{},
		},
	}
	routing.RegisterRoute(group, applicationRoutes, log)
}
