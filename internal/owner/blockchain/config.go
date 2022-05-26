package blockchain

import (
	"backend/tool"
	"encoding/json"
	"github.com/gochain/gochain/v3/accounts/abi"
	"github.com/gochain/web3"
	"github.com/spf13/viper"
)

const (
	name             = "BlockApp"
	symbol           = "BA"
	connContractAddr = "0x70833A8Db44Be18a80c9D5F698e952A5d93ED0ec"
)

var (
	privateKey   string
	nftContract  map[string]string
	connContract map[string]string
	channelUri   string
	abies        = map[string]abi.ABI{}
)

func init() {
	if nftContract == nil {
		nftContract = map[string]string{}
		tool.ReadConfig("./configs/owner/blockchain", "nft", "json")
		bytes, err := json.Marshal(viper.Get("abi"))
		if err != nil {
			panic(err)
		}
		nftContract["abi"] = string(bytes)
		nftContract["bytecode"] = viper.GetString("bytecode")
	}

	if connContract == nil {
		connContract = map[string]string{}
		tool.ReadConfig("./configs/owner/blockchain", "nft", "json")
		bytes, err := json.Marshal(viper.Get("abi"))
		if err != nil {
			panic(err)
		}
		connContract["abi"] = string(bytes)
	}

	if _, ok := abies["connection"]; !ok {
		tmp, err := web3.GetABI("./configs/owner/blockchain/connection.abi")
		if err != nil {
			panic(err)
		}
		abies["connection"] = *tmp
	}

	if privateKey == "" {
		tool.ReadConfig("./configs/owner/blockchain", "account", "yaml")
		privateKey = viper.GetString("privateKey")
	}

	if channelUri == "" {
		tool.ReadConfig("configs/owner", "uri", "yaml")
		channelUri = viper.GetString("channel")
	}
}
