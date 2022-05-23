package blockchain

import (
	"github.com/gochain/gochain/v3/common/hexutil"
	"github.com/gochain/gochain/v3/crypto"
)

const message = "Hello ArtBlock"
const messageLen = "14"
const klaytnMsg = "\x19Klaytn Signed Message:\n" + messageLen + message

func UnsignedAddress(signedAddr string) (string, error) {
	formattedMsg := []byte(klaytnMsg)
	sign, err := hexutil.Decode(signedAddr)
	if err != nil {
		return "", err
	}
	if sign[64] >= 2 {
		sign[64] = 1 - (sign[64] % 2)
	}
	hash := crypto.Keccak256(formattedMsg)
	pub, err := crypto.SigToPub(hash, sign)
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(*pub)
	return address.Hex(), nil
}
