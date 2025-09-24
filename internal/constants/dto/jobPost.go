package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type CreateJobPostRequest struct {
	ID          string    `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Deadline    time.Time `json:"deadline" binding:"required"`

	Responsibilities []string `json:"responsibilities" binding:"required,dive,required"`
	Requirements     []string `json:"requirements" binding:"required,dive,required"`

	DescriptionEmbedding      pgvector.Vector `json:"description_embedding"`
	RequirementsEmbedding     pgvector.Vector `json:"requirements_embedding"`
	ResponsibilitiesEmbedding pgvector.Vector `json:"responsibilities_embedding"`
}

type CreateJobPostResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Responsibilities []string  `json:"responsibilities"`
	Requirements     []string  `json:"requirements"`
	UserID           uuid.UUID `json:"user_id"`
	Deadline         time.Time `json:"deadline"`
	// CreatedAt   time.Time `json:"created_at"`
}

type GetAllJobPostsResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Responsibilities []string  `json:"responsibilities"`
	Requirements     []string  `json:"requirements"`
	UserID           uuid.UUID `json:"user_id"`
	Deadline         time.Time `json:"deadline"`
}

// UpdateJobPostRequest supports partial updates; omitted fields are ignored.
type UpdateJobPostRequest struct {
	ID               string     `json:"id"`
	Title            *string    `json:"title,omitempty"`
	Description      *string    `json:"description,omitempty"`
	Deadline         *time.Time `json:"deadline,omitempty"`
	Responsibilities *[]string  `json:"responsibilities,omitempty"`
	Requirements     *[]string  `json:"requirements,omitempty"`

	DescriptionEmbedding      *pgvector.Vector `json:"-"`
	RequirementsEmbedding     *pgvector.Vector `json:"-"`
	ResponsibilitiesEmbedding *pgvector.Vector `json:"-"`
}

type UpdateJobPostResponse struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Responsibilities []string  `json:"responsibilities"`
	Requirements     []string  `json:"requirements"`
	UserID           uuid.UUID `json:"user_id"`
	Deadline         time.Time `json:"deadline"`
}
