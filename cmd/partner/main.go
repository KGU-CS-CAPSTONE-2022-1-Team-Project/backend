package main

import (
	"backend/api/partner"
	pb "backend/proto/partner"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var port string

func init() {
	if !viper.IsSet("port") {
		tool.ReadConfig("./configs/partner", "client_secret", "json")
	}
	port = viper.GetString("port")
}

func main() {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterPartnerServiceServer(s, &partner.PartnerService{})
	reflection.Register(s)
	tool.Logger().Info("start GRPC partner server")
	if err := s.Serve(listen); err != nil {
		tool.Logger().Error("failed to serve", err, "port", port)
	}
}
