package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type CreateJobPostRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required,uuid"`
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
