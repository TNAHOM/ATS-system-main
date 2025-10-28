package initiator

import (
	"github.com/TNAHOM/ATS-system-main/internal/module"
	application "github.com/TNAHOM/ATS-system-main/internal/module/application"
	jobPost "github.com/TNAHOM/ATS-system-main/internal/module/jobPost"
	"github.com/TNAHOM/ATS-system-main/internal/module/user"
	"go.uber.org/zap"
)

type Module struct {
	User        module.User
	JobPost     module.JobPost
	Application module.Application
}

func InitModule(logger *zap.Logger, storage *Persistance, platform *Platform) *Module {
	return &Module{
		User:        user.Init(logger, storage.User, platform.Encryption),
		JobPost:     jobPost.Init(logger, storage.JobPost),
		Application: application.Init(logger, storage.Application),
	}
}
