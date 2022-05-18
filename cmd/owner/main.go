package main

import (
	"backend/api/owner/google"
	"backend/proto/owner/pb"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
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
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterOwnerServer(s, &google.OwnerService{})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
