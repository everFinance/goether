package example

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/everFinance/goether"
)

func TestContractCall(t *testing.T) {
	abi := `[{"constant": true,"inputs": [{"name": "","type": "address"}],"name": "balanceOf","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "dst","type": "address"},{"name": "wad","type": "uint256"}],"name": "transfer","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"}]`
	contractAddr := common.HexToAddress("0xd0a1e359811322d97991e03f863a0c30c2cf029c")
	prvHex := "your prvkey"
	rpc := "https://kovan.infura.io/v3/{{InfuraKey}}"

	testContract, err := goether.NewContract(contractAddr, abi, rpc, nil)
	if err != nil {
		panic(err)
	}

	amount, err := testContract.CallMethod("balanceOf", "latest", common.HexToAddress("0xa06b79e655db7d7c3b3e7b2cceeb068c3259d0c9"))
	t.Log(amount, err)

	testWallet, _ := goether.NewWallet(prvHex, rpc)
	testContract.Wallet = testWallet
	txHash, err := testContract.ExecMethod("transfer", nil, common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"), big.NewInt(100))
	t.Log(txHash, err)
}
