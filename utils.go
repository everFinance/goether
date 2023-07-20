package goether

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func EthToBN(amount float64) (bn *big.Int) {
	bf := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1000000000000000000))
	bn, _ = bf.Int(bn)
	return bn
}

func GweiToBN(amount float64) (bn *big.Int) {
	bf := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1000000000))
	bn, _ = bf.Int(bn)
	return bn
}

func EIP712Hash(typedData apitypes.TypedData) (hash []byte, err error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash = crypto.Keccak256(rawData)
	return
}

func Ecrecover(hash, signature []byte) (publicBy []byte, address common.Address, err error) {
	sig := make([]byte, len(signature))
	copy(sig, signature)
	if len(sig) != 65 {
		err = fmt.Errorf("invalid length of signture: %d", len(sig))
		return
	}

	if sig[64] != 27 && sig[64] != 28 && sig[64] != 1 && sig[64] != 0 {
		err = fmt.Errorf("invalid signature type")
		return
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	publicBy, err = crypto.Ecrecover(hash, sig)
	if err != nil {
		err = fmt.Errorf("can not ecrecover: %v", err)
		return
	}

	address = common.BytesToAddress(crypto.Keccak256(publicBy[1:])[12:])
	return
}

// Encrypt encrypt
func Encrypt(publicKey string, message []byte) ([]byte, error) {
	pub := common.FromHex(publicKey)
	pubKey, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return nil, err
	}
	eciesPub := ecies.ImportECDSAPublic(pubKey)
	result, err := ecies.Encrypt(rand.Reader, eciesPub, message, nil, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
