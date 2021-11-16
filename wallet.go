package goether

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/everFinance/ethrpc"
)

type TxOpts struct {
	Nonce     *int
	GasLimit  *int
	GasPrice  *big.Int
	GasTipCap *big.Int
	GasFeeCap *big.Int
}

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

	return NewWallet(strings.TrimSpace(string(b)), rpc)
}

func (w *Wallet) SendTx(to common.Address, amount *big.Int, data []byte, opts *TxOpts) (txHash string, err error) {
	opts, err = w.initTxOpts(to, amount, data, opts)
	if err != nil {
		return
	}

	if amount == nil {
		amount = big.NewInt(0)
	}

	tx, err := w.Signer.SignTx(
		*opts.Nonce, to, amount,
		*opts.GasLimit, opts.GasTipCap, opts.GasFeeCap,
		data, w.ChainID)
	if err != nil {
		return
	}

	raw, err := tx.MarshalBinary()
	if err != nil {
		return
	}

	return w.Client.EthSendRawTransaction(hexutil.Encode(raw))
}

func (w *Wallet) SendLegacyTx(to common.Address, amount *big.Int, data []byte, opts *TxOpts) (txHash string, err error) {
	opts, err = w.initTxOpts(to, amount, data, opts)
	if err != nil {
		return
	}

	if amount == nil {
		amount = big.NewInt(0)
	}
	tx, err := w.Signer.SignLegacyTx(
		*opts.Nonce, to, amount,
		*opts.GasLimit, opts.GasPrice,
		data, w.ChainID)
	if err != nil {
		return
	}

	raw, err := tx.MarshalBinary()
	if err != nil {
		return
	}

	return w.Client.EthSendRawTransaction(hexutil.Encode(raw))
}

func (w *Wallet) initTxOpts(to common.Address, amount *big.Int, data []byte, opts *TxOpts) (*TxOpts, error) {
	var (
		nonce, gasLimit int
		gasPrice        big.Int
		err             error
	)

	if opts == nil {
		opts = &TxOpts{}
	}

	if opts.Nonce == nil {
		nonce, err = w.GetPendingNonce()
		if err != nil {
			return nil, err
		}
		opts.Nonce = &nonce
	}

	if opts.GasLimit == nil {
		ethrpcTx := ethrpc.T{
			From:  w.Address.String(),
			To:    to.String(),
			Value: amount,
			Data:  hexutil.Encode(data),
		}
		gasLimit, err = w.Client.EthEstimateGas(ethrpcTx)
		if err != nil {
			return nil, err
		}
		opts.GasLimit = &gasLimit
	}

	if opts.GasPrice == nil {
		gasPrice, err = w.Client.EthGasPrice()
		if err != nil {
			return nil, err
		}
		opts.GasPrice = &gasPrice
	}

	if opts.GasTipCap == nil || opts.GasFeeCap == nil {
		opts.GasTipCap = opts.GasPrice
		opts.GasFeeCap = opts.GasPrice
	}

	return opts, nil
}

func (w *Wallet) GetAddress() string {
	return w.Address.String()
}

func (w *Wallet) GetNonce() (nonce int, err error) {
	return w.Client.EthGetTransactionCount(w.GetAddress(), "latest")
}

func (w *Wallet) GetPendingNonce() (nonce int, err error) {
	return w.Client.EthGetTransactionCount(w.GetAddress(), "pending")
}

func (w *Wallet) GetBalance() (balance big.Int, err error) {
	return w.Client.EthGetBalance(w.GetAddress(), "latest")
}
