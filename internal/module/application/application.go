package applicantmodule

import (
	"context"

	enum "github.com/TNAHOM/ATS-system-main/internal/constants/Enum"
	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/internal/module"
	"github.com/TNAHOM/ATS-system-main/internal/storage"
	"go.uber.org/zap"
)

type ApplicantModule struct {
	log              *zap.Logger
	applicantStorage storage.Application
}

func Init(log *zap.Logger, applicantStorage storage.Application) module.Application {
	return &ApplicantModule{log: log, applicantStorage: applicantStorage}
}

func (a *ApplicantModule) GetApplicationsByJobPostID(ctx context.Context, jobPostID string, progressStatus enum.ProgressStatus) ([]dto.GetApplicationsResponse, error) {
	applicants, err := a.applicantStorage.GetApplicationsByJobPostID(ctx, jobPostID, progressStatus)
	if err != nil {
		a.log.Error("Failed to get applications by job post ID", zap.Error(err))
		return nil, err
	}

	return applicants, nil
}
