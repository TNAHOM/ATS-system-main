package encryption

import (
	"errors"
	"os"
	"time"

	"github.com/TNAHOM/ATS-system-main/internal/constants/dto"
	"github.com/TNAHOM/ATS-system-main/platform"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SignedDetails struct {
	ID        string
	Id        string
	FirstName string
	LastName  string
	Email     string
	UserType  string
	jwt.StandardClaims
}

type TokenRelated struct {
	log *zap.Logger
}

func getSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func Init(log *zap.Logger) platform.Encryption {
	return &TokenRelated{log: log}
}

func (t *TokenRelated) GenerateToken(tokenField dto.GenerateUpdateToken) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     tokenField.Email,
		FirstName: tokenField.FirstName,
		LastName:  tokenField.LastName,
		ID:        tokenField.ID,
		UserType:  tokenField.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	PrivateKey := getSecretKey()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(PrivateKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(PrivateKey))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil

}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", errors.New("internal server error")
	}
	return string(hashed), nil
}

func VerifyPassword(hashedPass string, providedPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(providedPass))
	if err != nil {
		return false, err
	}

	return true, nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg error) {
	PrivateKey := getSecretKey()

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(PrivateKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("the token is invalid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
