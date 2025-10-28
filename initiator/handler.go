package initiator

import (
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	application "github.com/TNAHOM/ATS-system-main/internal/handler/application"
	jobpost "github.com/TNAHOM/ATS-system-main/internal/handler/jobPost"
	user "github.com/TNAHOM/ATS-system-main/internal/handler/user"
	"go.uber.org/zap"
)

type Handler struct {
	User        handler.User
	JobPost     handler.JobPost
	Application handler.Application
}

func InitHandler(log *zap.Logger, Module *Module) *Handler {
	return &Handler{
		User:        user.Init(log, Module.User),
		JobPost:     jobpost.Init(log, Module.JobPost),
		Application: application.Init(log, Module.Application),
	}
}
