package google

import (
	"backend/internal/owner/blockchain"
	"backend/tool"
	"context"
	"github.com/gochain/web3"
	"time"
)

func RegisterContract(address string, userId string) {
	pureClient, err := web3.Dial("https://api.baobab.klaytn.net:8651")
	defer pureClient.Close()
	if err != nil {
		tool.Logger().Warning("fail create web3 client", err)
		return
	}
	client := blockchain.WrappingClient{Client: pureClient}
	hash, err := client.DeployNFTMaker(address, userId)
	if err != nil {
		tool.Logger().Warning("fail DeployNFTMaker", err, "address", address, "uid", userId)
		return
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	contract, err := client.WaitDeploy(hash, timeout)
	if err != nil {
		tool.Logger().Warning("fail DeployNFTMaker", err, "ContractHash", hash.Hex())
		return
	}
	err = client.RegisterDeployedNFTMaker(address, contract)
	if err != nil {
		tool.Logger().Warning("fail DeployNFTMaker", err, "ContractHash", contract.Hex())
		return
	}
}
