package initiator

import (
	jobpost "github.com/TNAHOM/ATS-system-main/internal/glue/jobPost"
	"github.com/TNAHOM/ATS-system-main/internal/glue/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoute(grp *gin.RouterGroup, handler *Handler, module *Module, log *zap.Logger) {
	user.Init(grp, log, handler.User)
	jobpost.Init(grp, log, handler.JobPost)
}
