package test

import (
	"context"
	"encoding/json"
	"github.com/gochain/gochain/v3/common"
	"github.com/gochain/web3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

var testAddr string
var privateKey string
var connectionAddr string
var nftAbiString string
var bytecode string
var contractAddr common.Address

func init() {
	if !viper.IsSet("address") || !viper.IsSet("key") || viper.IsSet("connectionAddr") {
		viper.AddConfigPath("./configs/owner/blockchain")
		viper.SetConfigType("yaml")
		viper.SetConfigName("test_info")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	testAddr = viper.GetString("address")
	privateKey = viper.GetString("key")
	connectionAddr = viper.GetString("connectionAddr")
	if !viper.IsSet("bytecode") || !viper.IsSet("abi") {
		viper.SetConfigName("NFT")
		viper.SetConfigType("json")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
	byteArray, err := json.Marshal(viper.Get("abi"))
	if err != nil {
		panic(err)
	}
	nftAbiString = string(byteArray)
	bytecode = viper.GetString("bytecode")
}

func TestNFTContract(t *testing.T) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	client, err := web3.Dial("https://api.baobab.klaytn.net:8651")
	defer client.Close()

	chainId, err := client.GetChainID(timeout)
	require.Nilf(t, err, "체인아이디 조회중 에러발생\n\n %s", err)
	assert.Equal(t, chainId, big.NewInt(1001), "잘못된 체인아이디", "chainId: ", chainId)

	gasPrice, err := client.GetGasPrice(timeout)
	require.Nil(t, err, "gas price 조회중 에러발생", err)

	tx, err := web3.DeployContract(timeout,
		client,
		privateKey, bytecode, nftAbiString, gasPrice, 8_500_000,
		"name_test", "symbol_test", testAddr, connectionAddr, "test_url")
	require.Nil(t, err, "deployment 에러발생", err)
	var receipt *web3.Receipt
	for {
		receipt, err = client.GetTransactionReceipt(timeout, tx.Hash)
		if err == nil {
			break
		}
	}
	assert.NotEqual(t, receipt.ContractAddress.Hex(), "", "contract 생성 실패")
	contractAddr = receipt.ContractAddress
	t.Log(tx.Hash.Hex())
	t.Log(contractAddr.Hex())
}

func TestCallRegister(t *testing.T) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	myAbi, err := web3.GetABI("./configs/owner/blockchain/connection.abi")
	require.Nil(t, err, "abi 조회 실패", err)

	client, err := web3.Dial("https://api.baobab.klaytn.net:8651")
	defer client.Close()
	require.Nil(t, err, "클라이언트 연결 실패")
	gasPrice, err := client.GetGasPrice(timeout)
	require.Nil(t, err, "gas price 조회중 에러발생", err)

	t.Log(testAddr, contractAddr.Hex())
	result, err := web3.CallFunctionWithArgs(timeout,
		client, privateKey, connectionAddr, big.NewInt(0),
		gasPrice, 8_500_000,
		*myAbi, myAbi.Methods["register"].Name,
		testAddr, contractAddr.Hex(),
	)
	require.Nil(t, err, "function call 에러발생", err)
	t.Log(result.Hash.Hex())
}
