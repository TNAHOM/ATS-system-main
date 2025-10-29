package middleware

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware handler configured from environment variables.
// Behavior:
// - When ENV=production: reads CORS_ALLOWED_ORIGINS (comma-separated) and applies a stricter policy.
// - Otherwise: uses a permissive but practical set of localhost origins for development.
func CORS() gin.HandlerFunc {
	var corsConfig cors.Config
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		allowed := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowed == "" {
			// conservative fallback; update to your real production domain via env var
			allowed = "https://your-production-domain.com"
			log.Printf("CORS_ALLOWED_ORIGINS not set; falling back to %s", allowed)
		}
		corsConfig = cors.Config{
			AllowOrigins:     strings.Split(allowed, ","),
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
	} else {
		// local development defaults (tight but practical)
		corsConfig = cors.Config{
			AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080", "http://localhost:3000", "http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
	}

	return cors.New(corsConfig)
}
