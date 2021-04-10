package goether

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/signer/core"
	"github.com/stretchr/testify/assert"
)

func TestEthToBN(t *testing.T) {
	assert.Equal(t, big.NewInt(123456789876543200), EthToBN(0.1234567898765432))
}

func TestGweiToBN(t *testing.T) {
	assert.Equal(t, big.NewInt(1100000000), GweiToBN(1.1))
}

func TestEIP712(t *testing.T) {
	raw := `{"types": {"EIP712Domain": [{"name": "name","type": "string"},{"name": "version","type": "string"},{"name": "chainId","type": "uint256"}],"Order": [{"name": "action","type": "string"},{"name": "orderHashes","type": "string[]"},{"name": "makerAddress","type": "address"}]},"primaryType": "Order","domain": {"name": "ZooDex","version": "1","chainId": "42"},"message": {"action": "cancelOrder","orderHashes": ["0x123", "0x456", "0x789"],"makerAddress": "0xf9593A9d7F735814B87D08e8D8aD624f58d53B10"}}
	`
	typedData := core.TypedData{}
	err := json.Unmarshal([]byte(raw), &typedData)
	assert.NoError(t, err)

	hash, err := EIP712Hash(typedData)
	assert.NoError(t, err)
	assert.Equal(t, "0xcf3985dd9eb11ce656eafc2dddd08ce3058ad00c74669b3d171f31e9a0472d8e", hexutil.Encode(hash))

	addr, err := Ecrecover(hash, hexutil.MustDecode("0xa9a3e5f72b48651b735d0908f1f240b06eafe7166dbe6b4fc8b57d8b8515ef555fe4b124c2b50d6907423426ec46bc12c5956942dcfd01e02d70912c87a389c41b"))
	assert.NoError(t, err)
	assert.Equal(t, "0xf9593A9d7F735814B87D08e8D8aD624f58d53B10", addr.String())
}

func TestEcrecover(t *testing.T) {
	hash := accounts.TextHash([]byte("123"))
	sig, _ := hexutil.Decode("0x409c16579b4fc162f199f897497f5142101992af82cc6a0b9521413cf721151817e52781c0341fa333cdfea6ebe945b9231f8a8b3df7e7040203f9d7df26c2f21c")

	addr, err := Ecrecover(hash, sig)
	assert.NoError(t, err)
	assert.Equal(t, "0xab6c371B6c466BcF14d4003601951e5873dF2AcA", addr.String())
	// run again
	addr, err = Ecrecover(hash, sig)
	assert.NoError(t, err)
	assert.Equal(t, "0xab6c371B6c466BcF14d4003601951e5873dF2AcA", addr.String())
}
