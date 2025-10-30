package platform

import (
	"errors"
	"strings"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// ParsePagination binds query params into the simple dto.PaginationRequest you provided,
// applies defaults and clamps, and validates basic constraints.
func ParsePagination(c *gin.Context) (dto.PaginationRequest, error) {
	var p dto.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		return p, err
	}

	// normalize zero or negative values
	if p.Page <= 0 {
		p.Page = DefaultPage
	}
	if p.Size <= 0 {
		p.Size = DefaultPageSize
	}
	if p.Size > MaxPageSize {
		p.Size = MaxPageSize
	}

	// Optional: if you ever add sort/order fields, validate them here.
	// Example guard to prevent malicious values if you extend DTO later:
	_ = strings.ToLower // keep import used if you extend

	// final sanity check
	if p.Page <= 0 || p.Size <= 0 {
		return p, errors.New("invalid pagination parameters")
	}
	return p, nil
}

// Offset returns SQL offset for limit/offset queries
func Offset(p dto.PaginationRequest) int {
	return (p.Page - 1) * p.Size
}
