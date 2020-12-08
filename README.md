# goether

A simple Ethereum wallet implementation and utilities in Golang.

## Install

```shell
go get -u github.com/everFinance/goether
```

## Examples

```golang
prvHex := "your prvkey"
rpc := "https://infura.io/v3/{{InfuraKey}}"

testWallet, err := goether.NewWallet(prvHex, rpc)
if err != nil {
  panic(err)
}

txHash, err := testWallet.SendTx(nil, common.HexToAddress("0xa06b79E655Db7D7C3B3E7B2ccEEb068c3259d0C9"), goether.EthToBN(0.12), big.NewInt(100000), goether.GweiToBN(2), []byte("123"))
fmt.Println(txHash, err)
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

### Utils

- [x] EthToBN
- [x] GweiToBN
- [x] Ecrecover
