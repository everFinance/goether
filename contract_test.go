package goether

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestDecodeDataHex(t *testing.T) {
	abi := `[{"constant": true,"inputs": [{"name": "","type": "address"}],"name": "balanceOf","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "dst","type": "address"},{"name": "wad","type": "uint256"}],"name": "transfer","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"}]`
	testContract, err := NewContract(common.HexToAddress("0x0"), abi, "", nil)
	if err != nil {
		panic(err)
	}

	inputData, err := testContract.EncodeData(
		"transfer",
		common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"),
		big.NewInt(100))
	assert.NoError(t, err)

	methodName, params, err := testContract.DecodeDataHex(hexutil.Encode(inputData))
	assert.NoError(t, err)
	assert.Equal(t, "transfer", methodName)
	assert.Equal(t, common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"), params["dst"])
	assert.Equal(t, big.NewInt(100), params["wad"])

	_, _, err = testContract.DecodeDataHex("0xa9059c")
	assert.Equal(t, "data is too short", err.Error())
}

func TestDecodeEventHex(t *testing.T) {
	abi := `[{"anonymous": false,"inputs": [{"indexed": true,"name": "from","type": "address"},{"indexed": true,"name": "to","type": "address"},{"indexed": false,"name": "value","type": "uint256"}],"name": "Transfer","type": "event"}]`
	testContract, err := NewContract(common.HexToAddress("0x0"), abi, "", nil)
	if err != nil {
		panic(err)
	}

	eventName, values, err := testContract.DecodeEventHex(
		[]string{
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
			"0x000000000000000000000000a06b79e655db7d7c3b3e7b2cceeb068c3259d0c9",
			"0x0000000000000000000000003dd22a3ad30df8acaf12def3b27e085525a98065",
		},
		"0x0000000000000000000000000000000000000000000000000000000000989680",
	)
	assert.NoError(t, err)
	assert.Equal(t, "Transfer", eventName)
	assert.Equal(t, common.HexToAddress("0xa06b79e655db7d7c3b3e7b2cceeb068c3259d0c9"), values["from"])
	assert.Equal(t, common.HexToAddress("0x3dd22a3ad30df8acaf12def3b27e085525a98065"), values["to"])
	assert.Equal(t, big.NewInt(10000000), values["value"])
}
