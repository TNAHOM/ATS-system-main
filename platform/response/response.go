package response

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SendError standardizes error responses
func SendError(c *gin.Context, status int, message string, err interface{}) {
	// If the error is from validator, format it
	if verrs, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]string)
		for _, fe := range verrs {
			// fe.Tag() gives validation type ("required", "email", "min", etc.)
			details[fe.Field()] = fe.Tag()
		}
		c.AbortWithStatusJSON(status, dto.ErrorResponse{
			Code:    status,
			Message: message,
			Details: details,
		})
		return
	}

	// Default (non-validation errors)
	c.AbortWithStatusJSON(status, dto.ErrorResponse{
		Code:    status,
		Message: message,
		Details: err,
	})
}

// SendSuccess standardizes success responses
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
