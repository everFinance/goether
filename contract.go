package goether

import (
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/everFinance/ethrpc"
)

type Contract struct {
	Address common.Address
	ABI     abi.ABI

	Wallet *Wallet
	Client *ethrpc.EthRPC
}

func NewContract(address common.Address, abiStr, rpc string, wallet *Wallet) (*Contract, error) {
	Abi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}

	return &Contract{
		Address: address,
		ABI:     Abi,
		Wallet:  wallet,
		Client:  ethrpc.New(rpc),
	}, nil
}

// CallMethod Only read contract status
// tag:
// 	HEX String - an integer block number
//  String "earliest" for the earliest/genesis block
//  String "latest" - for the latest mined block
//  String "pending" - for the pending state/transactions
func (c *Contract) CallMethod(methodName, tag string, args ...interface{}) (res string, err error) {
	data, err := c.ABI.Pack(methodName, args...)
	if err != nil {
		return
	}

	return c.Client.EthCall(ethrpc.T{
		Data: hexutil.Encode(data),
		To:   c.Address.String(),
		From: c.Address.String(),
	}, tag)
}

// ExecMethod Execute tx
func (c *Contract) ExecMethod(methodName string, opts *TxOpts, args ...interface{}) (txHash string, err error) {
	if c.Wallet == nil {
		err = errors.New("wallet is nil")
		return
	}

	data, err := c.ABI.Pack(methodName, args...)
	if err != nil {
		return
	}

	return c.Wallet.SendTx(c.Address, big.NewInt(0), data, opts)
}

func (c *Contract) GetAddress() string {
	return c.Address.String()
}
