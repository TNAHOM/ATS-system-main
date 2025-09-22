package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

type JobPost struct {
	ID     string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`

	Responsibilities pq.StringArray `gorm:"type:text[]" json:"responsibilities"`
	Requirements     pq.StringArray `gorm:"type:text[]" json:"requirements"`

	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	DescriptionEmbedding      pgvector.Vector `gorm:"type:vector(3072)" json:"description_embedding"`
	RequirementsEmbedding     pgvector.Vector `gorm:"type:vector(3072)" json:"requirements_embedding"`
	ResponsibilitiesEmbedding pgvector.Vector `gorm:"type:vector(3072)" json:"responsibilities_embedding"`
}
