package goether

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var TestSigner *Signer

func init() {
	TestSigner, _ = NewSigner("8eda9cd543eaa0484b70e5dcf03ad23a65c01610e835cbef891bd7c59d965632")
}

func TestAddress(t *testing.T) {
	assert.Equal(t, "0xab6c371B6c466BcF14d4003601951e5873dF2AcA", TestSigner.Address.String())
}

func TestSignTx(t *testing.T) {
	_, err := TestSigner.SignTx(1, common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"), big.NewInt(0), 21000, big.NewInt(100000000000), nil, big.NewInt(42))
	assert.NoError(t, err)
}

func TestSignMsg(t *testing.T) {
	sign, err := TestSigner.SignMsg([]byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, "0x409c16579b4fc162f199f897497f5142101992af82cc6a0b9521413cf721151817e52781c0341fa333cdfea6ebe945b9231f8a8b3df7e7040203f9d7df26c2f21c", hexutil.Encode(sign))
}
