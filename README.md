# goether

A simple Ethereum wallet implementation and utilities in Golang.

## Install

```shell
go get -u github.com/everFinance/goether
```

## Examples

[More examples...](./example)

Send Eth with data
```golang
prvHex := "your prvkey"
rpc := "https://infura.io/v3/{{InfuraKey}}"

testWallet, err := goether.NewWallet(prvHex, rpc)
if err != nil {
  panic(err)
}

txHash, err := testWallet.SendTx(
  common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"), // To
  goether.EthToBN(0.12), // Value
  []byte("123"),         // Data
  nil)

// send with opts
nonce := int(1)
gasLimit := int(999999)
txHash, err := testWallet.SendTx(
  common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"),
  goether.EthToBN(0.12),
  []byte("123"),
  &goether.TxOpts{ // Configure nonce/gas yourself
    Nonce: &nonce,
    GasLimit: &gasLimit,
    GasPrice: goether.GweiToBN(10),
  })
```

Contract Interaction
```golang
abi := `[{"constant": true,"inputs": [{"name": "","type": "address"}],"name": "balanceOf","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "dst","type": "address"},{"name": "wad","type": "uint256"}],"name": "transfer","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"}]`
contractAddr := common.HexToAddress("contract address")
prvHex := "your prvkey"
rpc := "https://kovan.infura.io/v3/{{InfuraKey}}"

// init contract instance with wallet
testWallet, _ := goether.NewWallet(prvHex, rpc)
testContract, err := goether.NewContract(contractAddr, abi, rpc, testWallet)
if err != nil {
  panic(err)
}

// ERC20 BalanceOf
amount, err := testContract.CallMethod(
  "balanceOf", // Method name
  "latest", // Tag
  common.HexToAddress("0xa06b79e655db7d7c3b3e7b2cceeb068c3259d0c9")) // Args

// ERC20 Transfer
txHash, err := testContract.ExecMethod(
  "transfer", // Method name
  &goether.TxOpts{ // Configure nonce/gas yourself, nil load params from eth-node
    Nonce: &nonce,
    GasLimit: &gasLimit,
    GasPrice: goether.GweiToBN(10),
  },
  common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"), // Args
  big.NewInt(100))
```

## Modules

### Signer

Ethereum Account which can be used to sign messages and transactions.

- [x] SignTx
- [x] SignMsg
- [ ] SignTypedData

### Wallet

Connect to Ethereum Network, execute state changing operations.

- [x] SendTx
- [x] GetAddress
- [x] GetBalance
- [x] GetNonce
- [x] GetPendingNonce

### Contract

Creating Contract Instance for call & execute contract.

- [x] CallMethod
- [x] ExecMethod
- [x] EncodeData
- [x] EncodeDataHex
- [x] DecodeData
- [x] DecodeDataHex
- [x] DecodeEvent
- [x] DecodeEventHex

### Utils

- [x] EthToBN
- [x] GweiToBN
- [x] Ecrecover
