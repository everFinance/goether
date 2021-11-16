package example

import (
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

	// without opts
	txHash, err := testWallet.SendTx(
		common.HexToAddress("0xa2026731B31E4DFBa78314bDBfBFDC8cF5F761F8"), // To
		goether.EthToBN(0.12), // Value
		[]byte("123"),         // Data
		nil)
	t.Log(txHash, err)

	// set gasLimit & gasPrice
	gasLimit := int(999999)
	txHash, err = testWallet.SendTx(
		common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"),
		goether.EthToBN(0.12), []byte("123"),
		&goether.TxOpts{
			GasLimit: &gasLimit,
			GasPrice: goether.GweiToBN(2.1),
		})
	t.Log(txHash, err)
}
