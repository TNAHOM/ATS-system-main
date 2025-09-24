package jobpost

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
	jobPostHandler handler.JobPost,
) {
	jobPostRoutes := []routing.Route{
		{
			Method:  http.MethodPost,
			Path:    "/jobPost/create",
			Handler: jobPostHandler.CreateJobPost,
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
				// middleware.AuthUserTypeMiddleware(log, "recruiter"),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/jobPost/getAllJobPosts",
			Handler: jobPostHandler.GetAllJobPosts,
			Middleware: []gin.HandlerFunc{
				middleware.AuthMiddleware(log),
			},
		},
	}
	routing.RegisterRoute(group, jobPostRoutes, log)
}
