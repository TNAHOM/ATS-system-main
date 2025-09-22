package jobpost

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type jobPost struct {
	log           *zap.Logger
	jobPostModule module.JobPost
}

func Init(log *zap.Logger, jobPostModule module.JobPost) handler.JobPost {
	return &jobPost{log: log, jobPostModule: jobPostModule}
}

func (j *jobPost) CreateJobPost(ctx *gin.Context) {
	var jobPostModel dto.CreateJobPostRequest
	if err := ctx.ShouldBindJSON(&jobPostModel); err != nil {
		j.log.Error(err.Error(), zap.Any("request", jobPostModel))
		response.SendError(ctx, http.StatusBadRequest, "validation failed", err)
		return
	}

	res, err := j.jobPostModule.CreateJobPost(ctx, jobPostModel)
	if err != nil {
		j.log.Error(err.Error(), zap.Any("request", jobPostModel))
		response.SendError(ctx, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	ctx.JSON(http.StatusOK, dto.Envelope[dto.CreateJobPostResponse]{Data: res})
}
