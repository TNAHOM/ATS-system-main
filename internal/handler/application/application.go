package application

import (
	"net/http"

	enum "github.com/TNAHOM/ATS-system-main/internal/constants/Enum"
	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Application struct {
	log         *zap.Logger
	application module.Application
}

func Init(log *zap.Logger, application module.Application) handler.Application {
	return &Application{log, application}
}

// GetApplicationsByJobPostId godoc
// @Summary Get applications by job post ID
// @Tags applications
// @Accept json
// @Produce json
// @Param jobPostID path string true "Job Post ID"
// @Param progressStatus query string false "Progress Status" enums(APPLIED,INTERVIEWING,REJECTED,HIRED,SHORTLISTED)
// @Success 200 {object} dto.EnvelopeGetApplicationsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /applications/jobPost/{jobPostID} [get]
func (a *Application) GetApplicationsByJobPostID(ctx *gin.Context) {
	jobPostID := ctx.Param("jobPostID")
	getProgressStatus := ctx.DefaultQuery("progressStatus", string(enum.APPLIED))
	if !enum.IsValidProgressStatus(getProgressStatus) {
		a.log.Warn("Invalid progressStatus", zap.String("progressStatus", getProgressStatus))
		response.SendError(ctx, http.StatusBadRequest, "Invalid progressStatus", nil)
		return
	}
	progressStatus := enum.ProgressStatus(getProgressStatus)
	if err := uuid.Validate(jobPostID); err != nil {
		a.log.Warn("Invalid jobPostID format")
		response.SendError(ctx, http.StatusBadRequest, "Invalid jobPostID format", nil)
		return
	}

	applications, err := a.application.GetApplicationsByJobPostID(ctx, jobPostID, enum.ProgressStatus(progressStatus))
	if err != nil {
		a.log.Error(err.Error(), zap.String("jobPostID", jobPostID))
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	ctx.JSON(http.StatusOK, dto.EnvelopeGetApplicationsResponse{Data: applications})
}
