package google

import (
	"backend/internal/owner/blockchain"
	"context"
	"github.com/gochain/web3"
	"github.com/pkg/errors"
	"log"
	"time"
)

func RegisterContract(address string, userId string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	pureClient, err := web3.Dial("https://api.baobab.klaytn.net:8651")
	defer pureClient.Close()
	if err != nil {
		log.Println("fail create pure client", err)
		return
	}
	client := blockchain.WrappingClient{Client: pureClient}
	hash, err := client.DeployNFTMaker(address, userId)
	if err != nil {
		log.Println(errors.Cause(err))
		return
	}
	contract, err := client.WaitDeploy(hash, timeout)
	if err != nil {
		log.Println(errors.Cause(err))
		return
	}
	err = client.RegisterDeployedNFTMaker(address, contract)
	if err != nil {
		log.Println(errors.Cause(err))
		return
	}
}
