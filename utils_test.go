package goether

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
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

	tx := &types.Transaction{}
	// eip1559 tx
	data := []byte(`{
			"accessList": [],
			"blockHash": "0xe936ee0e5a915b9c163a7a1ff67269dd5f1ccb981f91b269a2130711e6a62598",
			"blockNumber": "0xa2f4f1",
			"chainId": "0x3",
			"from": "0xcba9f09e7e6b4a41a9d11f347416b75ee100344f",
			"gas": "0x5208",
			"gasPrice": "0x3b9aca0a",
			"hash": "0x87f6d244a9f99b2fb173f3eece4beecabb88da46c90f558b4d57fd10d8cb247b",
			"input": "0x",
			"maxFeePerGas": "0x22ecb25c00",
			"maxPriorityFeePerGas": "0x3b9aca00",
			"nonce": "0x1",
			"r": "0x6f741f168a787c708be5d7670de97b05c7e0a2e2f1959dd11fe5213508ea44a0",
			"s": "0x2e409f1b92769885c45f87bbdb858740e0bf7b74c1cd8297bcefbccb6b35ac21",
			"to": "0xc8d40bbc3ea2018ab8d76dbbee32f906207966d0",
			"transactionIndex": "0x72",
			"type": "0x2",
			"v": "0x0",
			"value": "0xaa87bee538000"
	}`)

	err = tx.UnmarshalJSON(data)
	assert.NoError(t, err)
	sigHash := types.NewLondonSigner(tx.ChainId()).Hash(tx)

	// assemble london tx signature
	sig = make([]byte, 0, 65)
	v, r, s := tx.RawSignatureValues()
	sig = append(r.Bytes(), s.Bytes()...)
	sig = append(sig, byte(v.Uint64()))
	// verify
	addr, err = Ecrecover(sigHash.Bytes(), sig)
	assert.NoError(t, err)
	assert.Equal(t, "0xcba9f09e7e6b4a41a9d11f347416b75ee100344f", strings.ToLower(addr.String()))
}
