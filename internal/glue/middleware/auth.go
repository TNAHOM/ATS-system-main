package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TNAHOM/ATS-system-main/platform/encryption"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Warn("Missing Authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn("invalid token format")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			ctx.Abort()
			return
		}

		trimmedToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := encryption.ValidateToken(trimmedToken)
		if err != nil {
			log.Error(err.Error(), zap.Any("request", map[string]interface{}{"authHeader": authHeader, "trimmed": trimmedToken}))
			ctx.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		// ctx := context.WithValue(ctx.Request.Context(), "claims", claims)
		// ctx.Request = ctx.Request.WithContext(ctx)
		fmt.Println("claims: s", claims)

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func AuthUserTypeMiddleware(log *zap.Logger, userType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		setUserType := ctx.GetString("user_type")

		if userType != setUserType {
			log.Error("userType not allowed", zap.Any("request", userType))
			ctx.JSON(http.StatusForbidden, gin.H{"error": "user type not allowed"})

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// Forward request to another microservice
func ProxyHandler(targetURL string, log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, exists := ctx.Get("claims")
		if !exists {
			log.Error("claims not found in context")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		fmt.Println("claims: s", v)

		claims, ok := v.(*encryption.SignedDetails)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx.Request.Context(), ctx.Request.Method, targetURL+ctx.Request.RequestURI, ctx.Request.Body)
		if err != nil {
			log.Error("gateway request failed", zap.Any("request", map[string]interface{}{
				"context": ctx.Request.Context(),
				"Method":  ctx.Request.Method,
				"url":     targetURL + ctx.Request.RequestURI,
				"body":    ctx.Request.Body,
			}))

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gateway request build failed"})
			return
		}
		// copy headers safely (preserve multiple values, avoid overwrite)
		for k, vv := range ctx.Request.Header {
			for _, hv := range vv {
				req.Header.Add(k, hv)
			}
		}
		req.Header.Set("X-User-Id", claims.ID)
		req.Header.Set("X-User-Email", claims.Email)
		req.Header.Set("X-User-First-Name", claims.FirstName)
		req.Header.Set("X-User-Last-Name", claims.LastName)
		req.Header.Set("X-User-Type", claims.UserType)

		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": "downstream service unreachable"})
			return
		}
		defer resp.Body.Close()
		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	}
}

// Forward request to another microservice
// func ProxyHandler(targetURL string) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// Take claims from context
// 		claims := ctx.Request.Context().Value("claims").(*encryption.SignedDetails)

// 		// Forward request to downstream service
// 		client := &http.Client{}
// 		req, err := http.NewRequest(ctx.Request.Method, targetURL+ctx.Request.RequestURI, ctx.Request.Body)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "gateway request build failed"})
// 			return
// 		}

// 		// Copy headers from original request
// 		for k, v := range ctx.Request.Header {
// 			req.Header[k] = v
// 		}

// 		// Optionally forward claims to downstream via custom header (don’t forward JWT directly if not needed)
// 		if claimsMap, ok := claims.(jwt.MapClaims); ok {
// 			for k, v := range claimsMap {
// 				req.Header.Set("X-Claim-"+k, fmt.Sprintf("%v", v))
// 			}
// 		}

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadGateway, gin.H{"error": "downstream service unreachable"})
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Relay response back to client
// 		ctx.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
// 	}
// }
