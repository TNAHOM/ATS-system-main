package applicantstorage

import (
	"context"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	log *zap.Logger
	db  *gorm.DB
}

func Init(log *zap.Logger, db *gorm.DB) storage.Application {
	return &Application{log, db}
}

func (a *Application) GetApplicationsByJobPostID(ctx context.Context, jobPostID string) ([]dto.GetApplicationsResponse, error) {
	var applications []dto.GetApplicationsResponse

	err := a.db.WithContext(ctx).Table("applications").Where("job_post_id = ? AND status = ?", jobPostID, "COMPLETED").Find(&applications).Error
	if err != nil {
		a.log.Error("Failed to get applications by job post ID", zap.Error(err))
		return nil, err
	}

	if len(applications) == 0 {
		a.log.Info("No applications found for the given job post ID", zap.String("jobPostID", jobPostID))
		return []dto.GetApplicationsResponse{}, nil
	}

	return applications, nil
}
