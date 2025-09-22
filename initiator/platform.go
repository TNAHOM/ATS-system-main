package initiator

import (
	"github.com/TNAHOM/ATS-system-main/platform"
	"github.com/TNAHOM/ATS-system-main/platform/encryption"
	"go.uber.org/zap"
)

type Platform struct {
	Encryption platform.Encryption
}

func InitPlatform(logger *zap.Logger) *Platform {
	return &Platform{
		Encryption: encryption.Init(logger),
	}
}
