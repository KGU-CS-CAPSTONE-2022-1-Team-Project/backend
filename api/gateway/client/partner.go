package client

import (
	partnerPb "backend/proto/partner"
	"backend/tool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var partnerClient partnerPb.PartnerServiceClient

var partnerOnce sync.Once

func Partner() partnerPb.PartnerServiceClient {
	if partnerClient == nil {
		partnerOnce.Do(func() {
			tool.ReadConfig("./config/gateway", "services", "yaml")
			info := viper.GetStringMapString("partner")
			conn, err := grpc.Dial(info["host"]+":"+info["port"], grpc.WithTransportCredentials(
				insecure.NewCredentials()),
			)
			if err != nil {
				panic(err)
			}
			partnerClient = partnerPb.NewPartnerServiceClient(conn)
		})
	}
	return partnerClient
}
