package client

import (
	ownerPb "backend/proto/owner"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var ownerClient ownerPb.OwnerClient

var ownerOnce sync.Once

func Owner() ownerPb.OwnerClient {
	if ownerClient == nil {
		ownerOnce.Do(func() {
			tool.ReadConfig("./config/gateway", "services", "yaml")
			info := viper.GetStringMapString("owner")
			conn, err := grpc.Dial(info["host"]+":"+info["port"], grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
			ownerClient = ownerPb.NewOwnerClient(conn)
		})
	}
	return ownerClient
}
