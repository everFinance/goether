package goether

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/everFinance/ethrpc"
)

type Wallet struct {
	Address common.Address
	ChainID *big.Int

	Signer *Signer
	Client *ethrpc.EthRPC
}

func NewWallet(prvHex, rpc string) (*Wallet, error) {
	signer, err := NewSigner(prvHex)
	if err != nil {
		return nil, err
	}

	client := ethrpc.New(rpc)

	version, err := client.NetVersion()
	if err != nil {
		return nil, err
	}
	chainID, ok := new(big.Int).SetString(version, 10)
	if !ok {
		return nil, fmt.Errorf("wrong chainID: %s", version)
	}

	return &Wallet{
		Address: signer.Address,
		ChainID: chainID,

		Signer: signer,
		Client: client,
	}, nil
}

func NewWalletFromPath(prvPath, rpc string) (*Wallet, error) {
	b, err := ioutil.ReadFile(prvPath)
	if err != nil {
		return nil, err
	}

	return NewWallet(string(b), rpc)
}

func (w *Wallet) SendTx(
	nonce *big.Int, to common.Address, amount *big.Int,
	gasLimit *big.Int, gasPrice *big.Int, data []byte,
) (txHash string, err error) {
	if nonce == nil {
		n, err := w.GetNonce()
		if err != nil {
			return "", err
		}

		nonce = big.NewInt(int64(n))
	}

	if amount == nil {
		amount = big.NewInt(0)
	}

	if gasLimit == nil {
		gas, err := w.Client.EthEstimateGas(ethrpc.T{
			From:  w.Address.String(),
			To:    to.String(),
			Value: amount,
			Data:  hexutil.Encode(data),
		})
		if err != nil {
			return "", err
		}

		gasLimit = big.NewInt(int64(gas))
	}

	tx, err := w.Signer.SignTx(
		int(nonce.Int64()), to, amount,
		int(gasLimit.Int64()), gasPrice, data, w.ChainID)
	if err != nil {
		return
	}

	raw, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return
	}

	return w.Client.EthSendRawTransaction(hexutil.Encode(raw))
}

func (w *Wallet) GetAddress() string {
	return w.Address.String()
}

func (w *Wallet) GetNonce() (nonce int, err error) {
	return w.Client.EthGetTransactionCount(w.GetAddress(), "latest")
}

func (w *Wallet) GetBalance() (balance big.Int, err error) {
	return w.Client.EthGetBalance(w.GetAddress(), "latest")
}
