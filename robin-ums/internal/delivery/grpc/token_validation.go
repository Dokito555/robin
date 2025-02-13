package grpc

import (
	"context"
	"fmt"

	"github.com/Dokito555/robin-ums/constants"
	token_validation_proto "github.com/Dokito555/robin-ums/internal/delivery/grpc/proto/token"
	"github.com/Dokito555/robin-ums/internal/services"
	"github.com/sirupsen/logrus"
)

type TokenValidationController struct {
	TokenService *services.TokenService
	Log *logrus.Logger
	token_validation_proto.UnimplementedTokenValidationServer
}

func NewTokenValidationController(service *services.TokenService, log *logrus.Logger) *TokenValidationController {
	return &TokenValidationController{
		Log: log,
		TokenService: service,
	}
}

func (s *TokenValidationController) ValidateToken(ctx context.Context, req *token_validation_proto.TokenRequest) (*token_validation_proto.TokenResponse, error) {
	var (
		token = req.GetToken()
	)

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

	return &token_validation_proto.TokenResponse{
		Message: constants.STATUS_SUCCESS,
		Data: &token_validation_proto.UserData{
			UserId: int64(claimToken.UserID),
			Email: claimToken.Email,
			Role: claimToken.Role,
		},
	}, nil
}