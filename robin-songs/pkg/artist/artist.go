package artist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	token_validation_proto "github.com/Dokito555/robin-songs/internal/delivery/grpc/proto/token"
	"github.com/Dokito555/robin-songs/utils/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Profile struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

type ArtistPkg struct {
	Log *logrus.Logger
	Config *viper.Viper
}

func NewArtist(logger *logrus.Logger, config *viper.Viper) *ArtistPkg {
	return &ArtistPkg{
		Log: logger,
		Config: config,
	}
}

func (s *ArtistPkg) GetArtist(ctx context.Context, token string) (*Profile, error) {
	url := s.Config.GetString("ARTIST_HOST")+s.Config.GetString("ARTIST_ENDPOINT_VERIFY")
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		s.Log.Warnf("failed to create http request to artist service: %v", err)
		return nil, errors.New("failed to create http request to artist service")
	}

	httpReq.Header.Set("Authorization", token)

	client := new(http.Client)
	rsp, err := client.Do(httpReq)
	if err != nil {
		s.Log.Warnf("failed to call artist service get profile: %+v", err)
		return nil, errors.New("failed to call artist service get profile")
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		s.Log.Warnf("failed to read response body: %v", err)
		return nil, errors.New("failed to read response body")
	}

	if rsp.StatusCode != http.StatusOK {
		s.Log.Warnf("got failed response from artist service: %d", rsp.StatusCode)
		return nil, fmt.Errorf("got failed response from artist service: %d", rsp.StatusCode)
	}

	pf := new(Profile)
	err = json.Unmarshal(body, pf)
	if err != nil {
		s.Log.Warnf("failed to decode response: %v", err)
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return pf, nil
}

func (s *ArtistPkg) ValidateToken(ctx context.Context, token string) (*Profile, error) {
	data := new(Profile)

	conn, err := grpc.Dial(
        s.Config.GetString("artist service_GRPC_HOST"), 
        grpc.WithInsecure(),
        grpc.WithBlock(),
        grpc.WithTimeout(5*time.Second),
    )

	if err != nil {
        return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
    }
    if conn == nil {
        return nil, fmt.Errorf("connection is nil after successful dial")
    }

	defer func() {
        if conn != nil {
            conn.Close()
        }
    }()

	client := token_validation_proto.NewTokenValidationClient(conn)

	req := &token_validation_proto.TokenRequest{
		Token: token,
	}

	response, err := client.ValidateToken(ctx, req)
	if err != nil {
        s.Log.Errorf("Token validation failed with error: %v", err)
        return data, fmt.Errorf("failed to validate token from grpc: %v", err)
    }

	if response == nil {
		return data, fmt.Errorf("response is nil")
	}

	if response.Message != constants.STATUS_SUCCESS {
		return data, fmt.Errorf("got response error from artist service: %s", response.Message)
	}

	if response.Data == nil {
        return data, fmt.Errorf("response data is nil")
    }

	data.ID = int(response.Data.UserId)
	data.Email = response.Data.Email
	data.UserName = response.Data.Username
	data.Role = response.Data.Role

	s.Log.Infof("Validated user profile: %+v", data)
	return data, nil
}