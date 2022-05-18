package client

import (
	"backend/proto/owner/pb"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var client pb.OwnerClient

var once sync.Once

func Owner() pb.OwnerClient {
	if client == nil {
		once.Do(func() {
			tool.ReadConfig("./config/gateway", "services", "yaml")
			info := viper.GetStringMapString("owner")
			conn, err := grpc.Dial(info["host"]+":"+info["port"], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
			client = pb.NewOwnerClient(conn)
		})
	}
	return client
}
