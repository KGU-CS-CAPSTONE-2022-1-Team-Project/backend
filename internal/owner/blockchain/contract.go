package blockchain

import (
	"context"
	"github.com/gochain/gochain/v3/common"
	"github.com/gochain/web3"
	"github.com/pkg/errors"
	"log"
	"math/big"
	"time"
)

type WrappingClient struct {
	web3.Client
}

func (receiver *WrappingClient) DeployNFTMaker(address, Id string) (common.Hash, error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	gasPrice, err := receiver.GetGasPrice(timeout)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "gas price error")
	}
	tx, err := web3.DeployContract(timeout, receiver, privateKey, nftContract["bytecode"],
		nftContract["abi"], gasPrice, 8_500_000,
		name, symbol, address, connContractAddr, channelUri+Id)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "Fail Deploy NFT")
	}
	return tx.Hash, nil
}

func (receiver *WrappingClient) WaitDeploy(hash common.Hash, ctx context.Context) (common.Address, error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	for {
		select {
		case <-ctx.Done():
			return common.Address{}, errors.Wrap(ctx.Err(), "timeout")
		default:
			receipt, err := receiver.GetTransactionReceipt(timeout, hash)
			if err == nil {
				return receipt.ContractAddress, nil
			}
		}
	}
}

func (receiver *WrappingClient) RegisterDeployedNFTMaker(address string, contract common.Address) error {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	gasPrice, _ := receiver.GetGasPrice(timeout)
	register, err := web3.CallTransactFunction(timeout,
		receiver, abies["connection"], address, privateKey,
		abies["connection"].Methods["register"].Name, big.NewInt(0), gasPrice, 8_500_000,
		address, contract)
	log.Println("registered:", register.Hash.Hex())
	if err != nil {
		return errors.Wrap(err, "Maker 등록 실패")
	}
	return nil
}
