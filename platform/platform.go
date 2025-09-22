package platform

import "github.com/TNAHOM/ATS-system-main/internal/constants/dto"

type Encryption interface {
	GenerateToken(tokenField dto.GenerateUpdateToken) (signedToken string, signedRefreshToken string, err error)
}
