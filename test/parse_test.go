package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	codes "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/yancaitech/go-polka/config"
	v11 "github.com/yancaitech/go-polka/model/v11"
	"github.com/yancaitech/go-polka/rpc"
	"github.com/yancaitech/go-polka/state"
	gu "github.com/yancaitech/go-utils"
)

func Test_ParseExtrinsic(t *testing.T) {
	/*
		client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
		//client, _ := rpc.New("wss://127.0.0.1", "", "")
		callIdx, err := client.GetCallIdx("Balances", "transfer")
		if err != nil {
			panic(err)
		}
		fmt.Println(callIdx)

		option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata}
	*/

	var option types.ScaleDecoderOption
	zbs, err := hex.DecodeString(config.OptionMetadata)
	if err != nil {
		panic(err)
	}
	optbs, err := gu.BytesUnlzma(zbs)
	err = json.Unmarshal(optbs, &option)
	if err != nil {
		panic(err)
	}

	extrinsic := "0x2d028440862cfb2c1560276fa4f8a28336fab90e80b7e57800d8317a872da6076bdb0f01c4cd77d09be2654f6658ae906d41fed12324319aff17a18af98f58c669097e00c13f7a20d92cb7ce8848f0e50567be7ac33d91158f9301c111477c7eac53a58b0008000500a075fdd61fc3088c07ea99717b3d7616cb09473e1601632a163c70b99378974702127a00"
	//extrinsic := "0x280403000b905e5b3f7501"
	//extrinsic := "0x39028440862cfb2c1560276fa4f8a28336fab90e80b7e57800d8317a872da6076bdb0f016043bbc6375286bc7f56390bfe52b20fa7aad0f56970010824da78c9c3bdf467179492060bc6ac02a183311f043985aa72987b6e8c7b3a925b6fb5c850f71080450000000500a075fdd61fc3088c07ea99717b3d7616cb09473e1601632a163c70b9937897470700c817a804"
	e := codes.ExtrinsicDecoder{}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsic)}, &option)
	e.Process()
	bb, err := json.MarshalIndent(e.Value, "", "    ")
	if err != nil {
		panic(err)
	}
	var resp v11.ExtrinsicDecodeResponse
	err = json.Unmarshal(bb, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bb))

	/*
		{
		    "account_id": "40862cfb2c1560276fa4f8a28336fab90e80b7e57800d8317a872da6076bdb0f",
		    "call_code": "0500",
		    "call_module": "Balances",
		    "call_module_function": "transfer",
		    "era": "00",
		    "extrinsic_hash": "1e33308814c9c602fbfdd03705bcbc224c22e416689388bf7f08f69a00d8bfc4",
		    "extrinsic_length": 141,
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
		            "value": "20000000000",
		            "value_raw": ""
		        }
		    ],
		    "signature": "befd1a8079b6b24015425b4c42c7ff16f2e329315e1851ff1cb244e51fc4540fe1c2ec7cd84705ada38550fce321e240c4e398fbcfa96d93ae0ffcee41a66886",
		    "tip": "0",
		    "version_info": "84"
		}
	*/
}
func Test_parseEvent(t *testing.T) {
	var (
		err  error
		key  string
		resp []byte
	)
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	blockHash := "0xcbdb8536a723abfedad1faa70e845da65b579260347e2681b64f7eff8619a0fe"
	key, err = state.CreateStorageKey(client.Metadata, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Rpc.SendRequest("state_getStorageAt", []interface{}{key, blockHash})
	if err != nil || len(resp) <= 0 {
		panic(err)
	}
	eventsHex := string(resp)
	//解析events
	option := types.ScaleDecoderOption{Metadata: &client.Metadata.Metadata, Spec: client.SpecVersion}
	ccHex := config.CoinEventType[client.CoinType]
	cc, _ := hex.DecodeString(ccHex)
	types.RegCustomTypes(source.LoadTypeRegistry(cc))
	e := codes.EventsDecoder{}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventsHex)}, &option)
	e.Process()
	data, err1 := json.Marshal(e.Value)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println(string(data))
}

func Test_GetCallIdx(t *testing.T) {
	//client, _ := rpc.New("http://ksm.rylink.io:30933", "", "")
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	callIdx, err := client.GetCallIdx("Balances", "transfer")
	if err != nil {
		panic(err)
	}
	fmt.Println(callIdx)
}
