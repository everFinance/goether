package goether

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
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

func EIP712Hash(typedData core.TypedData) (hash []byte, err error) {
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

func Ecrecover(hash, sig []byte) (addr common.Address, err error) {
	if len(sig) != 65 {
		err = fmt.Errorf("invalid length of signture: %d", len(sig))
		return
	}

	if sig[64] != 27 && sig[64] != 28 {
		err = fmt.Errorf("invalid signature type")
		return
	}
	sig[64] -= 27

	recoverPub, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		err = fmt.Errorf("can not ecrecover: %v", err)
		return
	}
	pubKey, err := crypto.UnmarshalPubkey(recoverPub)
	if err != nil {
		err = fmt.Errorf("can not unmarshal pubkey: %v", err)
		return
	}

	addr = crypto.PubkeyToAddress(*pubKey)
	return
}
