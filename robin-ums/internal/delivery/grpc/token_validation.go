package grpc

import (
	"context"
	"fmt"

	token_validation_proto "github.com/Dokito555/robin-ums/internal/delivery/grpc/proto/token"
	"github.com/Dokito555/robin-ums/internal/services"
	"github.com/Dokito555/robin-ums/utils/constants"
	"github.com/sirupsen/logrus"
)

// controller
type TokenValidationController struct {
	TokenService *services.TokenService
	Log          *logrus.Logger
	token_validation_proto.UnimplementedTokenValidationServer
}

func NewTokenValidationController(service *services.TokenService, log *logrus.Logger) *TokenValidationController {
	return &TokenValidationController{
		Log:          log,
		TokenService: service,
	}
}

func (s *TokenValidationController) ValidateToken(ctx context.Context, req *token_validation_proto.TokenRequest) (*token_validation_proto.TokenResponse, error) {
	var (
		token = req.GetToken()
	)

	s.Log.Info("Received token validation request")

	if token == "" {
		s.Log.Warnf("token is empty")
		err := fmt.Errorf("token is empty")
		return &token_validation_proto.TokenResponse{
			Message: err.Error(),
		}, err
	}

	claimToken, err := s.TokenService.ValidateToken(ctx, token)
	if err != nil {
		s.Log.Warnf("failed to validate token")
		return &token_validation_proto.TokenResponse{
			Message: err.Error(),
		}, err
	}

	if claimToken.UserID == 0 || claimToken.Role == "" {
        return nil, fmt.Errorf("invalid token data: missing required fields")
    }

	s.Log.WithFields(logrus.Fields{
        "userId": claimToken.UserID,
        "email":  claimToken.Email,
        "role":   claimToken.Role,
    }).Info("Token validated successfully")

	return &token_validation_proto.TokenResponse{
		Message: constants.STATUS_SUCCESS,
		Data: &token_validation_proto.UserData{
			UserId: int64(claimToken.UserID),
			Email:  claimToken.Email,
			Role:   claimToken.Role,
		},
	}, nil
}
