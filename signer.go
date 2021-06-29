package goether

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
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

	return NewSigner(strings.TrimSpace(string(b)))
}

func NewSignerFromMnemonic(mnemonic string) (*Signer, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0") // ethereum private path
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}
	priv, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	prvHex := crypto.FromECDSA(priv)
	return NewSigner(hex.EncodeToString(prvHex))
}

func (s Signer) GetPrivateKey() *ecdsa.PrivateKey {
	return s.key
}

func (s Signer) GetPublicKey() []byte {
	return crypto.FromECDSAPub(&s.key.PublicKey)
}

func (s Signer) GetPublicKeyHex() string {
	return hexutil.Encode(s.GetPublicKey())
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

func (s Signer) SignTypedData(typedData core.TypedData) (sig []byte, err error) {
	hash, err := EIP712Hash(typedData)
	if err != nil {
		return
	}

	sig, err = crypto.Sign(hash, s.key)
	if err != nil {
		return
	}

	sig[64] += 27
	return
}

// Decrypt decrypt
func (s Signer) Decrypt(ct []byte) ([]byte, error) {
	eciesPriv := ecies.ImportECDSA(s.key)
	return eciesPriv.Decrypt(ct, nil, nil)
}
