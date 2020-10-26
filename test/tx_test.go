package test

import (
	"fmt"
	"testing"

	"github.com/yancaitech/go-polka/rpc"
	"github.com/yancaitech/go-polka/tx"
)

func Test_BalanceTransfer(t *testing.T) {
	c, err := rpc.New("wss://rpc.polkadot.io", "", "")
	if err != nil {
		return
	}
	fmt.Println("1")
	btTx := tx.CreateTransaction("from", "to", 10000000, 12, 0)
	fmt.Println(btTx)
	btTx.SetGenesisHashAndBlockHash("genesisHash", "genesisHash", 0)
	// 通过方法去获取callIdx，不走config
	callIdx, err := c.GetCallIdx("Balances", "transfer")
	if err != nil {
		return
	}
	fmt.Println(callIdx)
	fmt.Println("2")
	btTx.SetSpecVersionAndCallId(uint32(c.SpecVersion), uint32(c.TransactionVersion), callIdx)
	_, message, err := btTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("3")
	sig, err := btTx.SignTransaction("private key", message)
	if err != nil {
		return
	}
	fmt.Println(sig)
	txHex, err := btTx.GetSignTransaction(sig)
	if err != nil {
		return
	}
	fmt.Println(txHex)
	//broadcast tx
	txidBytes, err := c.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txHex})
	if err != nil {
		return
	}
	txid := string(txidBytes)
	fmt.Println(txid)
}

func Test_UtilityBatch(t *testing.T) {
	c, err := rpc.New("wss://rpc.polkadot.io", "", "")
	if err != nil {
		return
	}
	address_amount := make(map[string]uint64)
	address_amount["to1"] = 123
	address_amount["to2"] = 456
	// .
	// .
	// .
	ubCallIdx, err := c.GetCallIdx("Utility", "batch")
	ubTx := tx.CreateUtilityBatchTransaction("from", ubCallIdx, 12, address_amount)
	ubTx.SetGenesisHashAndBlockHash("genesisHash", "genesisHash", 0)
	// 通过方法去获取callIdx，不走config
	callIdx, err := c.GetCallIdx("Balances", "transfer")
	if err != nil {
		return
	}
	ubTx.SetSpecVersionAndCallId(uint32(c.SpecVersion), uint32(c.TransactionVersion), callIdx)
	_, message, err := ubTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		return
	}
	sig, err := ubTx.SignTransaction("private key", message)
	if err != nil {
		return
	}
	txHex, err := ubTx.GetSignTransaction(sig)
	if err != nil {
		return
	}
	//broadcast tx
	txidBytes, err := c.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txHex})
	if err != nil {
		return
	}
	txid := string(txidBytes)
	fmt.Println(txid)
}
