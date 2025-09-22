package models

import "time"

type User struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName string `gorm:"column:first_name" json:"first_name" validate:"required"`
	LastName  string `gorm:"column:last_name" json:"last_name" validate:"required,min=3"`
	Password  string `gorm:"column:password" json:"password"`
	Email     string `gorm:"column:email" json:"email" validate:"required,email"`
	Phone     string `gorm:"column:phone" json:"phone"`
	UserType  string `gorm:"column:user_type" json:"user_type"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Token        string `gorm:"column:token" json:"token"`
	RefreshToken string `gorm:"column:refresh_token" json:"refresh_token"`

	// Relationship
	JobPosts []JobPost `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"job_posts"`
}
