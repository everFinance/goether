package goether

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer struct {
	Address common.Address
	key     *ecdsa.PrivateKey
}

func NewSigner(prvHex string) (*Signer, error) {
	k, err := crypto.HexToECDSA(prvHex)
	if err != nil {
		return nil, err
	}

	return &Signer{
		key:     k,
		Address: crypto.PubkeyToAddress(k.PublicKey),
	}, nil
}

func NewSignerFromPath(prvPath string) (*Signer, error) {
	b, err := ioutil.ReadFile(prvPath)
	if err != nil {
		return nil, err
	}

	return NewSigner(string(b))
}

func (s *Signer) SignTx(
	nonce int, to common.Address, amount *big.Int,
	gasLimit int, gasPrice *big.Int,
	data []byte, chainID *big.Int,
) (tx *types.Transaction, err error) {
	return types.SignTx(
		types.NewTransaction(
			uint64(nonce), to, amount,
			uint64(gasLimit), gasPrice, data),
		types.NewEIP155Signer(chainID),
		s.key,
	)
}

func (s Signer) SignMsg(msg []byte) (sig []byte, err error) {
	hash := accounts.TextHash(msg)
	sig, err = crypto.Sign(hash, s.key)
	if err != nil {
		return
	}

	sig[64] += 27
	return
}
