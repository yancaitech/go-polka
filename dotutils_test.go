package polkautils

import (
	"encoding/hex"
	"fmt"
	"hash"
	"testing"

	"github.com/dchest/blake2b"
	"github.com/yancaitech/go-polka/config"
	"github.com/yancaitech/go-polka/rpc"
	"github.com/yancaitech/go-polka/scale"
	"github.com/yancaitech/go-polka/tx"
	"github.com/yancaitech/go-polka/xxhash"
)

func TestGetStorage(t *testing.T) {
	addr := "12TbvYR6Jh65irwpUUGdLhdPCHP31dNZtrsT46iZR2Bqubob"
	ss, err := DotGetStorageFactor(addr)
	fmt.Println(ss, err)

	// balance example
	bal := "047374616b696e672092e97def2e010000000000000000000002"
	ns, err := scale.HexToBn(bal)
	fmt.Println(ns, err)

	var hasher hash.Hash = xxhash.New128(nil)
	hasher.Write([]byte("Balances"))
	bs := make([]byte, 0)
	bs = hasher.Sum(bs)
	ss = hex.EncodeToString(bs)
	fmt.Println(ss)
	// c2261276cc9d1f8598ea4b6a74b15c2f

	hasher.Reset()
	hasher.Write([]byte("FreeBalance"))
	bs = make([]byte, 0)
	bs = hasher.Sum(bs)
	ss = hex.EncodeToString(bs)
	fmt.Println(ss)
	// 6482b9ade7bc6657aaca787ba1add3b4

	pubk, err := DotPublicKeyFromAddress(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	bs, err = hex.DecodeString(pubk)
	if err != nil {
		fmt.Println(err)
		return
	}

	hasher = blake2b.New256()
	hasher.Write(bs)
	bs = make([]byte, 0)
	bs = hasher.Sum(bs)
	ss = hex.EncodeToString(bs)
	fmt.Println(ss)
}

func TestCreateAddress(t *testing.T) {
	secret := "789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c9"
	priv, err := hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pri:", secret)

	pubk, err := DotPublicKey(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	addr, err := DotAddress(priv, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println("Addr:", addr)

	pubk, err = DotPublicKeyFromAddress(addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	err = DotAddressValidate(addr, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
}

func TestCreateTransaction(t *testing.T) {
	secret := "789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c9"
	priv, err := hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pri:", secret)

	pubk, err := DotPublicKey(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	addr1, err := DotAddress(priv, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
	addr1 = "13Sr7g2gkwHyvgDN1RZUohnNXVD9SsPHp6r9VzFi97bqb43k"
	fmt.Println(addr1)

	secret2 := "789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c9"
	priv, err = hex.DecodeString(secret2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pri:", secret2)

	pubk, err = DotPublicKey(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	addr2, err := DotAddress(priv, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println(addr2)

	txid, txmsg, err := DotCreateSignedTransaction(addr2, addr1, 19809820000, 0, 0,
		uint32(25), uint32(5), "0500",
		"789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c9",
		"0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		"0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3", 2099908)
	if err != nil {
		panic(err)
	}
	fmt.Println(txmsg)

	txid, txjson, err := DotDecodeSignedTransaction(txmsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(txid, txjson)

	/*
		{
		    "account_id": "40862cfb2c1560276fa4f8a28336fab90e80b7e57800d8317a872da6076bdb0f",
		    "call_code": "0500",
		    "call_module": "Balances",
		    "call_module_function": "transfer",
		    "era": "4500",
		    "extrinsic_hash": "d5f9e5ea86b801bc0540cd36d8f1215f2f2796afdef21134b9213d4a95ed80bc",
		    "extrinsic_length": 140,
		    "nonce": 0,
		    "params": [
		        {
		            "name": "dest",
		            "type": "Address",
		            "value": "a075fdd61fc3088c07ea99717b3d7616cb09473e1601632a163c70b993789747",
		            "value_raw": ""
		        },
		        {
		            "name": "value",
		            "type": "Compact\u003cBalance\u003e",
		            "value": "1000000",
		            "value_raw": ""
		        }
		    ],
		    "signature": "98a2780c429457719024b76216d93816dc76a4f5a7e7ae469488055cab22602bc3150030cdd33a3d0b10a2a44c0bb81ad51581fce51ebdb929bea8cc3bcaa38e",
		    "tip": "1",
		    "version_info": "84"
		}
	*/

	/*
		client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
		//client, _ := rpc.New("wss://127.0.0.1", "", "")
		extrinsic := txmsg
		//extrinsic := "0x39028492e0feb85e225ee7c1800966ab1e69d2e0afe8021304421f9ca9bf5ea9f9784601b418c283bdadb2243dedb254d46af218c82e51048600432d473c2d0c6621d57860867599803edcc62c2bb57967dbb14b44c096bae7125142e8518be9e17e418df50200000500704f74890129ae1e780a1dcaa27fd395bfce5744d4e377dd197174537df36702070010a5d4e8"
		e := codes.ExtrinsicDecoder{}
		option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata}

		optbs, err := json.Marshal(option)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(optbs, &option)
		if err != nil {
			panic(err)
		}
		lzbs, err := utils.BytesLzma(optbs)
		if err != nil {
			panic(err)
		}
		optstr := hex.EncodeToString(lzbs)
		//fmt.Println(optstr)

		fileName := "e:/temp/optmeta.dat"
		dstFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer dstFile.Close()
		dstFile.WriteString(optstr)

		//fmt.Println("OPT:", option.Metadata)
		e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsic)}, &option)
		e.Process()
		bb, err := json.Marshal(e.Value)
		if err != nil {
			panic(err)
		}
		var resp v11.ExtrinsicDecodeResponse
		err = json.Unmarshal(bb, &resp)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bb))
	*/
}

func TestCreateTransaction2(t *testing.T) {
	secret := "789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c9"
	priv, err := hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pri:", secret)

	pubk, err := DotPublicKey(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	addr1, err := DotAddress(priv, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println(addr1)

	secret = "789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c8"
	priv, err = hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pri:", secret)

	pubk, err = DotPublicKey(priv)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pub:", pubk)

	addr2, err := DotAddress(priv, config.PolkadotPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println(addr2)

	txid, txmsg, err := DotCreateSignedTransaction(addr2, addr1, 20000000000, 0, 0, 18, 4, "0500", secret,
		"0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		"0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3", 2087748)
	if err != nil {
		panic(err)
	}
	fmt.Println(txmsg)

	txid, txjson, err := DotDecodeSignedTransaction(txmsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(txid, txjson)

	// client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	// //client, _ := rpc.New("wss://127.0.0.1", "", "")
	// //broadcast tx
	// txidBytes, err := client.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txmsg})
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }
	// txid := string(txidBytes)
	// fmt.Println("txid:", txid)
}

func TestTransfer(t *testing.T) {
	c, err := rpc.New("wss://rpc.polkadot.io", "", "")
	if err != nil {
		return
	}
	btTx := tx.CreateTransaction("12TbvYR6Jh65irwpUUGdLhdPCHP31dNZtrsT46iZR2Bqubob",
		"14dPiKPm6aGWwNWyZge96wdyrfrzWkuCPJDhm5TMCTPYJVKD", 2000000, 2, 0)
	btTx.SetGenesisHashAndBlockHash("0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		"0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3", 2287748)
	// 通过方法去获取callIdx，不走config
	callIdx, err := c.GetCallIdx("Balances", "transfer")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("callIdx:", callIdx, "specVer:", c.SpecVersion, "txVer:", c.TransactionVersion)

	btTx.SetSpecVersionAndCallId(uint32(c.SpecVersion), uint32(c.TransactionVersion), callIdx)
	_, message, err := btTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("MSG:", message)
	sig, err := btTx.SignTransaction("789108e387659c239c9ffcba7d022fad7e70460b2c7b6ad449a87c5fa450b6c8", message)
	if err != nil {
		fmt.Println(err)
		return
	}
	txHex, err := btTx.GetSignTransaction(sig)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(txHex)

	return

	//broadcast tx
	txidBytes, err := c.Rpc.SendRequest("author_submitExtrinsic", []interface{}{txHex})
	if err != nil {
		fmt.Println(err)
		return
	}
	txid := string(txidBytes)
	fmt.Println(txid)
}
