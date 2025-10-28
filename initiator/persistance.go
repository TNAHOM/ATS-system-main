package initiator

import (
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	application "github.com/TNAHOM/ATS-system-main/internal/storage/application"
	jobPost "github.com/TNAHOM/ATS-system-main/internal/storage/jobPost"
	"github.com/TNAHOM/ATS-system-main/internal/storage/user"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Persistance struct {
	User        storage.Users
	JobPost     storage.JobPosts
	Application storage.Application
}

func InitPersistance(db *gorm.DB, log *zap.Logger) *Persistance {
	return &Persistance{
		User:        user.Init(log, db),
		JobPost:     jobPost.Init(log, db),
		Application: application.Init(log, db),
	}
}
