package grpc

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)


type GrpcConfig struct {
	Log *logrus.Logger
	Viper *viper.Viper
	TokenValidationController *TokenValidationController
}

func (c *GrpcConfig) Setup() {
	c.RunGrpc()
}

func (c *GrpcConfig) RunGrpc() {
	lis, err := net.Listen("tcp", ":"+c.Viper.GetString("GRPC_PORT"))
	if err != nil {
		c.Log.Fatal("failed to listen to grpc port:", err)
	}

	s := grpc.NewServer()

	token_validation_proto.RegisterTokenValidationServer(s, c.TokenValidationController)
	
	logrus.Info("listening to grpc port: " + c.Viper.GetString("GRPC_PORT"))
	if err := s.Serve(lis); err != nil {
		c.Log.Fatal("failed to serve grpc port: ", err)
	}
}