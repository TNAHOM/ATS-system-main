package applicantstorage

import (
	"context"

	enum "github.com/TNAHOM/ATS-system-main/internal/constants/Enum"
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

func (a *Application) GetApplicationsByJobPostID(ctx context.Context, jobPostID string, progressStatus enum.ProgressStatus) (dto.GetApplicationsResponseWithMetaData, error) {
	var applications []dto.GetApplicationsResponse
	var metaData dto.GetMetaDataApplicationsResponse

	err := a.db.WithContext(ctx).Table("applications").Where("job_post_id = ? AND status = ? AND progress_status = ?", jobPostID, "COMPLETED", progressStatus).Find(&applications).Error
	if err != nil {
		a.log.Error("Failed to get applications by job post ID", zap.Error(err))
		return dto.GetApplicationsResponseWithMetaData{}, err
	}

	progressStatuses := []enum.ProgressStatus{
		enum.APPLIED,
		enum.SHORTLISTED,
		enum.INTERVIEWING,
		enum.HIRED,
		enum.REJECTED,
	}

	for _, status := range progressStatuses {
		var count int64
		err := a.db.WithContext(ctx).Table("applications").Where("job_post_id = ? AND status = ? AND progress_status = ?", jobPostID, "COMPLETED", status).Count(&count).Error
		if err != nil {
			a.log.Error("Failed to count applications by progress status", zap.Error(err))
			return dto.GetApplicationsResponseWithMetaData{}, err
		}
		switch status {
		case enum.APPLIED:
			metaData.AppliedCount = int(count)
		case enum.SHORTLISTED:
			metaData.ShortlistedCount = int(count)
		case enum.INTERVIEWING:
			metaData.InterviewingCount = int(count)
		case enum.HIRED:
			metaData.HiredCount = int(count)
		case enum.REJECTED:
			metaData.RejectedCount = int(count)
		}
	}

	if len(applications) == 0 {
		a.log.Info("No applications found for the given job post ID", zap.String("jobPostID", jobPostID))
		return dto.GetApplicationsResponseWithMetaData{Applications: []dto.GetApplicationsResponse{}, MetaData: metaData}, nil
	}

	return dto.GetApplicationsResponseWithMetaData{
		Applications: applications,
		MetaData:     metaData,
	}, nil
}

func (a *Application) UpdateApplicationProgressStatus(ctx context.Context, applicationID string, progressStatus enum.ProgressStatus) (dto.UpdateApplicationResponse, error) {
	result := a.db.WithContext(ctx).Table("applications").Where("id = ?", applicationID).Update("progress_status", progressStatus)
	if result.Error != nil {
		a.log.Error("Failed to update application progress status", zap.Error(result.Error))
		return dto.UpdateApplicationResponse{}, result.Error
	}

	var updatedApplication dto.UpdateApplicationResponse
	err := a.db.WithContext(ctx).Table("applications").Where("id = ?", applicationID).First(&updatedApplication).Error
	if err != nil {
		a.log.Error("Failed to fetch updated application", zap.Error(err))
		return dto.UpdateApplicationResponse{}, err
	}

	return updatedApplication, nil
}
