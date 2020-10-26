package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yancaitech/go-polka/rpc"
)

// state_getKeys
// twox_128("System")                 + twox_128("Account")
// 0x26aa394eea5630e07c48ae0c9558cef7 + 0xb99d880ec681799c0cf30e8886371da9
// 0x26aa394eea5630e07c48ae0c9558cef7b99d880ec681799c0cf30e8886371da9

// state_getStorage
// xxhash128("ModuleName") + xxhash128("StorageName") + blake256hash("FirstKey") + blake256hash("SecondKey")

/*
twox_128("Balances)                                              = "0xc2261276cc9d1f8598ea4b6a74b15c2f"
twox_128("FreeBalance")                                          = "0x6482b9ade7bc6657aaca787ba1add3b4"
scale_encode("5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY") = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"

blake2_128_concat("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d") = "0xde1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"

state_getStorage("0xc2261276cc9d1f8598ea4b6a74b15c2f6482b9ade7bc6657aaca787ba1add3b4de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d") = "0x0000a0dec5adc9353600000000000000"


twox_128("Balances)                                       = "0xc2261276cc9d1f8598ea4b6a74b15c2f"
twox_128("FreeBalance")                                   = "0x6482b9ade7bc6657aaca787ba1add3b4"

state_getKeys("0xc2261276cc9d1f8598ea4b6a74b15c2f6482b9ade7bc6657aaca787ba1add3b4") = [
    "0xc2261276cc9d1f8598ea4b6a74b15c2f6482b9ade7bc6657aaca787ba1add3b4de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d",
    "0xc2261276cc9d1f8598ea4b6a74b15c2f6482b9ade7bc6657aaca787ba1add3b432a5935f6edc617ae178fef9eb1e211fbe5ddb1579b72e84524fc29e78609e3caf42e85aa118ebfe0b0ad404b5bdd25f",
    ...
]
*/

func Test_GetBlock(t *testing.T) {
	client, _ := rpc.New("wss://rpc.polkadot.io", "", "")
	data, err := client.GetBlockByNumber(1801866)
	if err != nil {
		panic(err)
	}
	d, _ := json.Marshal(data)
	fmt.Println(string(d))
}
