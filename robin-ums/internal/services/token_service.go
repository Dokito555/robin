package services

import (
	"context"
	"fmt"
	"time"

	model "github.com/Dokito555/robin-ums/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenService struct {
    Log   *logrus.Logger
    Viper *viper.Viper
}

func NewTokenService(log *logrus.Logger, viper *viper.Viper) *TokenService {
	return &TokenService{
		Log: log,
		Viper: viper,
	}
}

var mapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 24,
	"refresh_token": time.Hour * 72,
}

func (s *TokenService) GenerateToken(ctx context.Context, userID int, tokenType string, email string, role string) (string, error) {
	_, exists := mapTypeToken[tokenType]
    if !exists {
        return "", fmt.Errorf("invalid token type: %s", tokenType)
    }

	claimToken := model.ClaimToken{
		UserID: userID,
		Email:  email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.Viper.GetString("APP_NAME"),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	resultToken, err := token.SignedString([]byte(s.Viper.GetString("APP_SECRET")))
	if err != nil {
		s.Log.Warnf("failed to generate token: %+v", err)
		return resultToken, fmt.Errorf("failed to generate token: %v", err)
	}
	return resultToken, err
}

func (s *TokenService) ValidateToken(ctx context.Context, token string) (*model.ClaimToken, error) {
	var (
		claimToken *model.ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &model.ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			s.Log.Warnf("failed to validate method jwt: %+v", t.Header["alg"])
			return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
		}
		return []byte(s.Viper.GetString("APP_SECRET")), nil
	})

	if err != nil {
		s.Log.Warnf("failed to parse jwt: %+v", err)
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*model.ClaimToken); !ok || !jwtToken.Valid {
		s.Log.Warnf("token invalid")
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}