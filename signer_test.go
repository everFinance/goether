package goether

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/signer/core"
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
	sig, err := TestSigner.SignMsg([]byte("123"))
	assert.NoError(t, err)
	assert.Equal(t, "0x409c16579b4fc162f199f897497f5142101992af82cc6a0b9521413cf721151817e52781c0341fa333cdfea6ebe945b9231f8a8b3df7e7040203f9d7df26c2f21c", hexutil.Encode(sig))
}

func TestSignTypedData(t *testing.T) {
	typedDataJson := `{"primaryType":"Mail","types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Person":[{"name":"name","type":"string"},{"name":"wallet","type":"address"}],"Mail":[{"name":"from","type":"Person"},{"name":"to","type":"Person"},{"name":"contents","type":"string"}]},"domain":{"name":"Ether Mail","version":"1","chainId":1,"verifyingContract":"0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC"},"message":{"from":{"name":"Cow","wallet":"0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},"to":{"name":"Bob","wallet":"0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},"contents":"Hello, Bob!"}}`
	var typedData core.TypedData
	json.Unmarshal([]byte(typedDataJson), &typedData)
	sig, err := TestSigner.SignTypedData(typedData)
	assert.NoError(t, err)
	assert.Equal(t, "0x0a80cc322f7a5e5e0965ff84fd76ea6479fae9d8ce29f7e076c3b9cb8e8097b80052ee91f1fad54d4348c518e42395848431746c8fee40420041f87ea05b5a281c", hexutil.Encode(sig))

	hash, err := EIP712Hash(typedData)
	addr, err := Ecrecover(hash, sig)
	assert.Equal(t, "0xab6c371B6c466BcF14d4003601951e5873dF2AcA", addr.String())
}

func TestSigner_Decrypt_Encrypt(t *testing.T) {
	signer, err := NewSigner("dde30fa25128addf45656a39c0570fd06fce3e48056457b9f1f9fda603cc4be1")
	assert.NoError(t, err)

	msg := "aaa bbb ccc ddd"
	pubkey := signer.GetPublicKeyHex()
	ct, err := Encrypt(pubkey, []byte(msg))
	assert.NoError(t, err)

	decMsg, err := signer.Decrypt(ct)
	assert.NoError(t, err)

	assert.Equal(t, msg, string(decMsg))

	// test02
	msg02 := "aaa 1234 你好 ..."
	pubkey = TestSigner.GetPublicKeyHex()
	ct, err = Encrypt(pubkey, []byte(msg02))
	assert.NoError(t, err)

	decMsg, err = TestSigner.Decrypt(ct)
	assert.NoError(t, err)

	assert.Equal(t, msg02, string(decMsg))
}
