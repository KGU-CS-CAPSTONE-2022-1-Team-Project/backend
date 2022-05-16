package test

import (
	"context"
	"encoding/json"
	"fmt"
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
var abiString string
var bytecode string

func init() {
	if !viper.IsSet("address") || !viper.IsSet("key") || viper.IsSet("connectionAddr") {
		viper.AddConfigPath("./configs/blockchain")
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
	abiString = string(byteArray)
	bytecode = viper.GetString("bytecode")
	fmt.Println("test")
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
		privateKey, bytecode, abiString, gasPrice, 8_500_000,
		"name_test", "symbol_test", testAddr, connectionAddr, "test_url")
	require.Nil(t, err, "deployment 에러발생", err)
	hex := tx.Hash.Hex()
	t.Log(hex)
}
