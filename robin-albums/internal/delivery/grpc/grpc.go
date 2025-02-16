package grpc

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)


type GrpcConfig struct {
	Log *logrus.Logger
	Viper *viper.Viper
}

func (c *GrpcConfig) Setup() {
	c.RunGrpc()
}

func (c *GrpcConfig) RunGrpc() {
	port := c.Viper.GetString("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
        c.Log.Fatalf("failed to listen on port %s: %v", port, err)
    }

	s := grpc.NewServer(
        grpc.KeepaliveParams(keepalive.ServerParameters{
            MaxConnectionIdle: 5 * time.Minute,
            Time:             20 * time.Second,
            Timeout:         10 * time.Second,
        }),
    )
	
	c.Log.Infof("gRPC server listening on port %s", port)
    if err := s.Serve(lis); err != nil {
        c.Log.Fatalf("Failed to serve gRPC: %v", err)
    }
}