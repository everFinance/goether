package example

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/everFinance/goether"
)

func TestWalletSend(t *testing.T) {
	prvHex := "your prvkey"
	rpc := "https://kovan.infura.io/v3/{{InfuraKey}}"

	testWallet, err := goether.NewWallet(prvHex, rpc)
	if err != nil {
		panic(err)
	}

	txHash, err := testWallet.SendTx(nil, common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"), goether.EthToBN(0.12), big.NewInt(100000), goether.GweiToBN(2), []byte("123"))
	fmt.Println(txHash, err)
}
