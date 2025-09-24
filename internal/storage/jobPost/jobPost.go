package jobpost

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	models "github.com/TNAHOM/ATS-system-main/internal/constants/model"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobPost struct {
	log *zap.Logger
	db  *gorm.DB
}

func Init(log *zap.Logger, db *gorm.DB) storage.JobPosts {
	return &JobPost{log, db}
}

func (j *JobPost) CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error) {
	params := models.JobPost{
		ID:               jobPost.ID,
		Title:            jobPost.Title,
		Description:      jobPost.Description,
		Responsibilities: jobPost.Responsibilities,
		Requirements:     jobPost.Requirements,

		DescriptionEmbedding:      jobPost.DescriptionEmbedding,
		RequirementsEmbedding:     jobPost.RequirementsEmbedding,
		ResponsibilitiesEmbedding: jobPost.ResponsibilitiesEmbedding,

		UserID:   jobPost.UserID,
		Deadline: jobPost.Deadline,
	}
	if err := j.db.Create(&params).Error; err != nil {
		j.log.Error("Failed to create job post", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	return dto.CreateJobPostResponse{
		ID:               jobPost.ID,
		Title:            jobPost.Title,
		Description:      jobPost.Description,
		Responsibilities: jobPost.Responsibilities,
		Requirements:     jobPost.Requirements,
		UserID:           jobPost.UserID,
		Deadline:         jobPost.Deadline,
	}, nil
}

func (j *JobPost) GetAllJobPosts(ctx context.Context) ([]dto.GetAllJobPostsResponse, error) {
	var jobPosts []models.JobPost
	if err := j.db.WithContext(ctx).Find(&jobPosts).Error; err != nil {
		j.log.Error("failed to fetch job posts", zap.Error(err))
		return nil, err
	}

	res := make([]dto.GetAllJobPostsResponse, len(jobPosts))
	for i, jp := range jobPosts {
		res[i] = dto.GetAllJobPostsResponse{
			ID:               jp.ID,
			Title:            jp.Title,
			Description:      jp.Description,
			Responsibilities: jp.Responsibilities,
			Requirements:     jp.Requirements,
			UserID:           jp.UserID,
			Deadline:         jp.Deadline,
		}
	}
	return res, nil
}

func (j *JobPost) GetJobPostByID(ctx context.Context, id string) (dto.UpdateJobPostResponse, error) {
	var jp models.JobPost
	if err := j.db.WithContext(ctx).First(&jp, "id = ?", id).Error; err != nil {
		j.log.Error("failed to get job post by id", zap.String("id", id), zap.Error(err))
		return dto.UpdateJobPostResponse{}, err
	}
	return dto.UpdateJobPostResponse{
		ID:               jp.ID,
		Title:            jp.Title,
		Description:      jp.Description,
		Responsibilities: jp.Responsibilities,
		Requirements:     jp.Requirements,
		UserID:           jp.UserID,
		Deadline:         jp.Deadline,
	}, nil
}

func (j *JobPost) UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.UpdateJobPostResponse, error) {
	// Build map of fields to update
	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Deadline != nil {
		updates["deadline"] = *req.Deadline
	}
	if req.Responsibilities != nil {
		updates["responsibilities"] = pq.StringArray(*req.Responsibilities)
	}
	if req.Requirements != nil {
		updates["requirements"] = pq.StringArray(*req.Requirements)
	}
	// embeddings if regenerated
	if req.DescriptionEmbedding != nil {
		updates["description_embedding"] = *req.DescriptionEmbedding
	}
	if req.RequirementsEmbedding != nil {
		updates["requirements_embedding"] = *req.RequirementsEmbedding
	}
	if req.ResponsibilitiesEmbedding != nil {
		updates["responsibilities_embedding"] = *req.ResponsibilitiesEmbedding
	}

	if len(updates) == 0 {
		// nothing to update, just return current state
		return j.GetJobPostByID(ctx, req.ID)
	}

	if err := j.db.WithContext(ctx).Model(&models.JobPost{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		j.log.Error("failed updating job post", zap.String("id", req.ID), zap.Error(err))
		return dto.UpdateJobPostResponse{}, err
	}

	return j.GetJobPostByID(ctx, req.ID)
}
