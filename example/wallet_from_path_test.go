package example

import (
	"testing"

	"github.com/everFinance/goether"
)

func TestWalletFromPath(t *testing.T) {
	path := "./prvkey"
	rpc := "https://kovan.infura.io/v3/{{InfuraKey}}"

	testWallet, err := goether.NewWalletFromPath(path, rpc)
	if err != nil {
		panic(err)
	}

	t.Log(testWallet.Address.String())
}
