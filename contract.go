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
	data, err := c.EncodeData(methodName, args...)
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

	data, err := c.EncodeData(methodName, args...)
	if err != nil {
		return
	}

	return c.Wallet.SendTx(c.Address, big.NewInt(0), data, opts)
}

func (c *Contract) GetAddress() string {
	return c.Address.String()
}

func (c *Contract) EncodeData(methodName string, args ...interface{}) ([]byte, error) {
	return c.ABI.Pack(methodName, args...)
}

func (c *Contract) EncodeDataHex(methodName string, args ...interface{}) (hex string, err error) {
	by, err := c.EncodeData(methodName, args...)
	if err != nil {
		return
	}

	return hexutil.Encode(by), nil
}

func (c *Contract) DecodeData(data []byte) (methodName string, params map[string]interface{}, err error) {
	if len(data) < 4 {
		err = errors.New("data is too short")
		return
	}

	method, err := c.ABI.MethodById(data[:4])
	if err != nil {
		return
	}
	methodName = method.Name

	params = make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(params, data[4:])
	return
}

func (c *Contract) DecodeDataHex(dataHex string) (methodName string, params map[string]interface{}, err error) {
	data := common.FromHex(dataHex)
	return c.DecodeData(data)
}

func (c *Contract) DecodeEvent(topics []common.Hash, data []byte) (eventName string, values map[string]interface{}, err error) {
	if len(topics) < 1 {
		err = errors.New("no topics found")
		return
	}

	event, err := c.ABI.EventByID(topics[0])
	if err != nil {
		return
	}
	eventName = event.Name

	values = make(map[string]interface{})
	// parse topics
	var indexed abi.Arguments
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	indexTopics := []common.Hash{}
	for _, topic := range topics[1:] {
		indexTopics = append(indexTopics, topic)
	}
	if err = abi.ParseTopicsIntoMap(values, indexed, indexTopics); err != nil {
		return
	}

	// parse data
	err = event.Inputs.UnpackIntoMap(values, data)
	return
}

func (c *Contract) DecodeEventHex(topicsHex []string, dataHex string) (eventName string, values map[string]interface{}, err error) {

	topics := []common.Hash{}
	for _, topicHex := range topicsHex {
		topics = append(topics, common.HexToHash(topicHex))
	}
	return c.DecodeEvent(topics, common.FromHex(dataHex))
}
