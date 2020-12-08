package goether

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
