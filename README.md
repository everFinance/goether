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
fmt.Println(txHash, err)
```

Integrate Contract
```golang
abi := `[{"constant": true,"inputs": [{"name": "","type": "address"}],"name": "balanceOf","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "dst","type": "address"},{"name": "wad","type": "uint256"}],"name": "transfer","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"}]`
contractAddr := common.HexToAddress("contract address")
prvHex := "your prvkey"
rpc := "https://kovan.infura.io/v3/{{InfuraKey}}"

testWallet, _ := goether.NewWallet(prvHex, rpc)
testContract, err := goether.NewContract(contractAddr, abi, rpc, testWallet)
if err != nil {
  panic(err)
}

// ERC20 BalanceOf
amount, err := testContract.CallMethod("balanceOf", "latest", common.HexToAddress("0xa06b79e655db7d7c3b3e7b2cceeb068c3259d0c9"))
// ERC20 Tansfer
txHash, err := testContract.ExecMethod("transfer", nil, common.HexToAddress("0xab6c371B6c466BcF14d4003601951e5873dF2AcA"), big.NewInt(100))
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

### Contract

### Utils

- [x] EthToBN
- [x] GweiToBN
- [x] Ecrecover
