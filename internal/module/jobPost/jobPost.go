package jobpost

import (
	"context"
	"fmt"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"github.com/TNAHOM/ATS-system-main/platform/ai"
	"github.com/TNAHOM/ATS-system-main/platform/encryption"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JobPost struct {
	log            *zap.Logger
	jobPostStorage storage.JobPosts
}

func Init(logger *zap.Logger, jobPostStorage storage.JobPosts) module.JobPost {
	return &JobPost{log: logger, jobPostStorage: jobPostStorage}
}

func (j *JobPost) CreateJobPost(ctx context.Context, jobPost dto.CreateJobPostRequest) (dto.CreateJobPostResponse, error) {
	client, err := ai.NewClient(ctx)
	if err != nil {
		j.log.Error("Failed to create AI client", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	claimsVal := ctx.Value("claims")
	if claimsVal == nil {
		j.log.Error("claims not found in context")
		return dto.CreateJobPostResponse{}, fmt.Errorf("invalid token")
	}

	userClaims, ok := claimsVal.(*encryption.SignedDetails)
	if !ok {
		j.log.Error("claims type assertion failed")
		return dto.CreateJobPostResponse{}, fmt.Errorf("invalid token")
	}

	jobPost.ID = uuid.New().String()
	uid, err := uuid.Parse(userClaims.ID)
	if err != nil {
		j.log.Error("Failed to parse userClaims.ID to uuid.UUID", zap.Error(err))
		return dto.CreateJobPostResponse{}, fmt.Errorf("invalid user ID format: %w", err)
	}
	jobPost.UserID = uid

	descriptionEmbedding, err := ai.Embedding(ctx, client, []string{jobPost.Description})
	if err != nil {
		j.log.Error("Failed to generate description embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	responsibilityEmbedding, err := ai.Embedding(ctx, client, jobPost.Responsibilities)
	if err != nil {
		j.log.Error("Failed to generate responsibility embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}
	requirementEmbedding, err := ai.Embedding(ctx, client, jobPost.Requirements)
	if err != nil {
		j.log.Error("Failed to generate requirement embedding", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	jobPost.DescriptionEmbedding = descriptionEmbedding
	jobPost.RequirementsEmbedding = requirementEmbedding
	jobPost.ResponsibilitiesEmbedding = responsibilityEmbedding

	res, err := j.jobPostStorage.CreateJobPost(ctx, jobPost)

	if err != nil {
		j.log.Error("Failed to create job post", zap.Error(err))
		return dto.CreateJobPostResponse{}, err
	}

	return dto.CreateJobPostResponse{
		ID:               res.ID,
		Title:            res.Title,
		Description:      res.Description,
		Responsibilities: res.Responsibilities,
		Requirements:     res.Requirements,
		UserID:           res.UserID,
		Deadline:         res.Deadline,
	}, nil
}

// func GetAllJobPosts() {}
func (j *JobPost) GetAllJobPosts(ctx context.Context) ([]dto.GetAllJobPostsResponse, error) {
	jobPosts, err := j.jobPostStorage.GetAllJobPosts(ctx)
	if err != nil {
		j.log.Error("failed to get all job posts", zap.Error(err))
		return nil, err
	}
	return jobPosts, nil
}

func (j *JobPost) UpdateJobPost(ctx context.Context, req dto.UpdateJobPostRequest) (dto.UpdateJobPostResponse, error) {
	// Get current record
	existing, err := j.jobPostStorage.GetJobPostByID(ctx, req.ID)
	if err != nil {
		return dto.UpdateJobPostResponse{}, err
	}

	// Determine which textual fields changed to decide on embedding regeneration
	needDescEmbed := req.Description != nil && *req.Description != existing.Description
	needRespEmbed := req.Responsibilities != nil
	if needRespEmbed && len(*req.Responsibilities) == len(existing.Responsibilities) {
		// shallow compare slices
		same := true
		for i, v := range *req.Responsibilities {
			if v != existing.Responsibilities[i] {
				same = false
				break
			}
		}
		if same {
			needRespEmbed = false
		}
	}
	needReqEmbed := req.Requirements != nil
	if needReqEmbed && len(*req.Requirements) == len(existing.Requirements) {
		same := true
		for i, v := range *req.Requirements {
			if v != existing.Requirements[i] {
				same = false
				break
			}
		}
		if same {
			needReqEmbed = false
		}
	}

	if needDescEmbed || needRespEmbed || needReqEmbed {
		client, err := ai.NewClient(ctx)
		if err != nil {
			j.log.Error("failed to init ai client", zap.Error(err))
			return dto.UpdateJobPostResponse{}, err
		}
		if needDescEmbed {
			emb, err := ai.Embedding(ctx, client, []string{*req.Description})
			if err != nil {
				return dto.UpdateJobPostResponse{}, err
			}
			req.DescriptionEmbedding = &emb
		}
		if needRespEmbed {
			emb, err := ai.Embedding(ctx, client, *req.Responsibilities)
			if err != nil {
				return dto.UpdateJobPostResponse{}, err
			}
			req.ResponsibilitiesEmbedding = &emb
		}
		if needReqEmbed {
			emb, err := ai.Embedding(ctx, client, *req.Requirements)
			if err != nil {
				return dto.UpdateJobPostResponse{}, err
			}
			req.RequirementsEmbedding = &emb
		}
	}

	updated, err := j.jobPostStorage.UpdateJobPost(ctx, req)
	if err != nil {
		return dto.UpdateJobPostResponse{}, err
	}
	return updated, nil
}

func (j *JobPost) DeleteJobPost(ctx context.Context, id string) error {
	return j.jobPostStorage.DeleteJobPost(ctx, id)
}
