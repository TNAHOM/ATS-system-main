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

func (j *JobPost) GetAllJobPostsByUserId(ctx context.Context, p dto.PaginationRequest) ([]dto.GetJobPostsResponse, int, error) {
	var jobPosts []models.JobPost
	userId, err := ctx.Value("UserID").(string)
	if !err {
		j.log.Error("UserID not found in context")
		return nil, 0, gorm.ErrInvalidData
	}

	base := j.db.WithContext(ctx).Model(&models.JobPost{}).Where("user_id = ?", userId)

	var total64 int64
	if err := base.Count(&total64).Error; err != nil {
		j.log.Error("failed to count job posts", zap.Error(err))
		return nil, 0, err
	}

	total := int(total64)

	offset := (p.Page - 1) * p.Size
	if err := base.Limit(p.Size).Offset(offset).Find(&jobPosts).Error; err != nil {
		j.log.Error("failed to fetch job posts", zap.Error(err))
		return nil, 0, err
	}

	res := make([]dto.GetJobPostsResponse, len(jobPosts))
	for i, jp := range jobPosts {
		res[i] = dto.GetJobPostsResponse{
			ID:               jp.ID,
			Title:            jp.Title,
			Description:      jp.Description,
			Responsibilities: jp.Responsibilities,
			Requirements:     jp.Requirements,
			UserID:           jp.UserID,
			Deadline:         jp.Deadline,
			CreatedAt:        jp.CreatedAt,
			UpdatedAt:        jp.UpdatedAt,
			ApplicantCount:   jp.ApplicantCount,
		}
	}
	if len(res) == 0 {
		return []dto.GetJobPostsResponse{}, total, nil
	}
	return res, total, nil
}

func (j *JobPost) GetAllJobPosts(ctx context.Context, p dto.PaginationRequest) ([]dto.GetJobPostsResponse, int, error) {
	var jobPosts []models.JobPost

	base := j.db.WithContext(ctx).Model(&models.JobPost{})

	var total64 int64
	if err := base.Count(&total64).Error; err != nil {
		j.log.Error("failed to count job posts", zap.Error(err))
		return nil, 0, err
	}

	total := int(total64)

	offset := (p.Page - 1) * p.Size
	if err := base.Limit(p.Size).Offset(offset).Find(&jobPosts).Error; err != nil {
		j.log.Error("failed to fetch job posts", zap.Error(err))
		return nil, 0, err
	}

	res := make([]dto.GetJobPostsResponse, len(jobPosts))
	for i, jp := range jobPosts {
		res[i] = dto.GetJobPostsResponse{
			ID:               jp.ID,
			Title:            jp.Title,
			Description:      jp.Description,
			Responsibilities: jp.Responsibilities,
			Requirements:     jp.Requirements,
			UserID:           jp.UserID,
			Deadline:         jp.Deadline,
			CreatedAt:        jp.CreatedAt,
			UpdatedAt:        jp.UpdatedAt,
			ApplicantCount:   jp.ApplicantCount,
		}
	}
	return res, total, nil
}

func (j *JobPost) GetJobPostByID(ctx context.Context, id string) (dto.GetJobPostsResponse, error) {
	var jp models.JobPost
	if err := j.db.WithContext(ctx).First(&jp, "id = ?", id).Error; err != nil {
		j.log.Error("failed to get job post by id", zap.String("id", id), zap.Error(err))
		return dto.GetJobPostsResponse{}, err
	}
	return dto.GetJobPostsResponse{
		ID:               jp.ID,
		Title:            jp.Title,
		Description:      jp.Description,
		Responsibilities: jp.Responsibilities,
		Requirements:     jp.Requirements,
		UserID:           jp.UserID,
		Deadline:         jp.Deadline,
		CreatedAt:        jp.CreatedAt,
		UpdatedAt:        jp.UpdatedAt,
		ApplicantCount:   jp.ApplicantCount,
	}, nil
}

func (j *JobPost) UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.GetJobPostsResponse, error) {
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
		return dto.GetJobPostsResponse{}, err
	}

	return j.GetJobPostByID(ctx, req.ID)
}

func (j *JobPost) DeleteJobPost(ctx context.Context, id string) error {
	tx := j.db.WithContext(ctx).Where("id = ?", id).Delete(&models.JobPost{})
	if tx.Error != nil {
		j.log.Error("failed to soft delete job post", zap.String("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
