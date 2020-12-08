package goether

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestEthToBN(t *testing.T) {
	assert.Equal(t, big.NewInt(123456789876543200), EthToBN(0.1234567898765432))
}

func TestGweiToBN(t *testing.T) {
	assert.Equal(t, big.NewInt(1100000000), GweiToBN(1.1))
}

func TestEcrecover(t *testing.T) {
	hash := accounts.TextHash([]byte("123"))
	sign, _ := hexutil.Decode("0x409c16579b4fc162f199f897497f5142101992af82cc6a0b9521413cf721151817e52781c0341fa333cdfea6ebe945b9231f8a8b3df7e7040203f9d7df26c2f21c")

	addr, err := Ecrecover(hash, sign)
	assert.NoError(t, err)
	assert.Equal(t, "0xab6c371B6c466BcF14d4003601951e5873dF2AcA", addr.String())
}
