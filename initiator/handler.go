package initiator

import (
	"github.com/TNAHOM/ATS-system-main/internal/handler"
	jobpost "github.com/TNAHOM/ATS-system-main/internal/handler/jobPost"
	"github.com/TNAHOM/ATS-system-main/internal/handler/user"
	"go.uber.org/zap"
)

type Handler struct {
	User    handler.User
	JobPost handler.JobPost
}

func InitHandler(log *zap.Logger, userModule *Module) *Handler {
	return &Handler{
		User:    user.Init(log, userModule.User),
		JobPost: jobpost.Init(log, userModule.JobPost),
	}
}
