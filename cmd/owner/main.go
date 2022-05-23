package main

import (
	"backend/api/owner/google"
	pb "backend/proto/owner"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

var port string

func init() {
	if !viper.IsSet("port") {
		tool.ReadConfig("./configs/owner", "client_secert", "json")
	}
	port = viper.GetString("port")
}

func main() {
	defer func() {
		r := recover()
		tool.Logger().Error("panic", r.(error))
		main()
	}()
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterOwnerServer(s, &google.OwnerService{})
	reflection.Register(s)
	tool.Logger().Info("start GRPC owner server", "port", port)
	if err := s.Serve(listen); err != nil {
		tool.Logger().Error("failed to serve", err, "port", port)
	}
}
