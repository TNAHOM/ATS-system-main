package jobpost

import (
	"net/http"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/platform"
	"github.com/TNAHOM/ATS-system-main/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type jobPost struct {
	log           *zap.Logger
	jobPostModule module.JobPost
}

func Init(log *zap.Logger, jobPostModule module.JobPost) handler.JobPost {
	return &jobPost{log: log, jobPostModule: jobPostModule}
}

// CreateJobPost godoc
// @Summary Create a job post
// @Tags JobPost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body dto.CreateJobPostRequest true "Create job post payload"
// @Success 200 {object} dto.EnvelopeCreateJobPostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/create [post]
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

// GetJobPostByID godoc
// @Summary Get a job post by ID
// @Tags JobPost
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job post ID"
// @Success 200 {object} dto.EnvelopeGetJobPostsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/getJobPostByID/{id} [get]
func (j *jobPost) GetJobPostByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendError(ctx, http.StatusBadRequest, "missing id", nil)
		return
	}
	jobPost, err := j.jobPostModule.GetJobPostByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SendError(ctx, http.StatusNotFound, "job post not found", nil)
			return
		}
		j.log.Error("failed to get job post by id", zap.Error(err))
		response.SendError(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	ctx.JSON(http.StatusOK, dto.Envelope[dto.GetJobPostsResponse]{Data: jobPost})
}

// GetAllJobPostsByUserId godoc
// @Summary List all job posts
// @Tags JobPost
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedGetJobPostsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/getAllJobPostsByUserId [get]
func (j *jobPost) GetAllJobPostsByUserId(ctx *gin.Context) {
	p, err := platform.ParsePagination(ctx)
	if err != nil {
		j.log.Error("invalid pagination params", zap.Error(err))
		response.SendError(ctx, http.StatusBadRequest, "invalid pagination params", err)
		return
	}

	items, meta, err := j.jobPostModule.GetAllJobPostsByUserId(ctx, p)
	if err != nil {
		j.log.Error("failed to get job posts", zap.Error(err))
		response.SendError(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	res := dto.PaginatedResponse[dto.GetJobPostsResponse]{
		Items: items,
		Meta:  meta,
	}
	ctx.JSON(http.StatusOK, dto.Envelope[dto.PaginatedResponse[dto.GetJobPostsResponse]]{Data: res})
}

// GetAllJobPosts godoc
// @Summary List all job posts
// @Tags JobPost
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedGetJobPostsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/getAllJobPosts [get]
func (j *jobPost) GetAllJobPosts(ctx *gin.Context) {
	p, err := platform.ParsePagination(ctx)
	if err != nil {
		j.log.Error("invalid pagination params", zap.Error(err))
		response.SendError(ctx, http.StatusBadRequest, "invalid pagination params", err)
		return
	}

	jobPosts, meta, err := j.jobPostModule.GetAllJobPosts(ctx, p)
	if err != nil {
		j.log.Error("failed to get job posts", zap.Error(err))
		response.SendError(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	res := dto.PaginatedResponse[dto.GetJobPostsResponse]{
		Items: jobPosts,
		Meta:  meta,
	}

	ctx.JSON(http.StatusOK, dto.Envelope[dto.PaginatedResponse[dto.GetJobPostsResponse]]{Data: res})
}

// UpdateJobPost godoc
// @Summary Update a job post (partial)
// @Tags JobPost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job post ID"
// @Param data body dto.UpdateJobPostRequest true "Update job post payload"
// @Success 200 {object} dto.EnvelopeUpdateJobPostResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/update/{id} [patch]
func (j *jobPost) UpdateJobPost(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendError(ctx, http.StatusBadRequest, "missing id", nil)
		return
	}
	var req dto.UpdateJobPostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		j.log.Error("invalid update payload", zap.Error(err))
		response.SendError(ctx, http.StatusBadRequest, "validation failed", err)
		return
	}
	if req.Title == nil && req.Description == nil && req.Deadline == nil && req.Responsibilities == nil && req.Requirements == nil {
		response.SendError(ctx, http.StatusBadRequest, "no fields provided to update", nil)
		return
	}
	req.ID = id
	updated, err := j.jobPostModule.UpdateJobPost(ctx, req)
	if err != nil {
		j.log.Error("failed to update job post", zap.Error(err))
		response.SendError(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	ctx.JSON(http.StatusOK, dto.Envelope[dto.GetJobPostsResponse]{Data: updated})
}

// DeleteJobPost godoc
// @Summary Delete a job post
// @Tags JobPost
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job post ID"
// @Success 200 {object} dto.EnvelopeAny
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /jobPost/{id} [delete]
func (j *jobPost) DeleteJobPost(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.SendError(ctx, http.StatusBadRequest, "missing id", nil)
		return
	}
	if err := j.jobPostModule.DeleteJobPost(ctx, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			response.SendError(ctx, http.StatusNotFound, "job post not found", nil)
			return
		}
		j.log.Error("failed to delete job post", zap.Error(err))
		response.SendError(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	ctx.JSON(http.StatusOK, dto.Envelope[any]{Data: gin.H{"deleted": true}})
}
