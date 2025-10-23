package initiator

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/TNAHOM/ATS-system-main/docs"
)

func InitSwagger(server *gin.RouterGroup) {
	docs.SwaggerInfo.Title = "ATS System API"
	docs.SwaggerInfo.Description = "REST API for ATS-system-main (users, auth, job posts)"
	docs.SwaggerInfo.Version = "1.0"

	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		swaggerHost = fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	}
	docs.SwaggerInfo.Host = swaggerHost

	// Ensure base path matches the API grouping used in the server
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Register the swagger handler under the provided group so the final
	// URL becomes /api/swagger/*any when the group is mounted at /api
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
