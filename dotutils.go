package polkautils

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/dchest/blake2b"
	codes "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	iutiles "github.com/itering/scale.go/utiles"
	"github.com/yancaitech/go-polka/config"
	v11 "github.com/yancaitech/go-polka/model/v11"
	"github.com/yancaitech/go-polka/ss58"
	"github.com/yancaitech/go-polka/tx"
	gu "github.com/yancaitech/go-utils"
)

// DotAddress func
func DotAddress(prikey, prefix []byte) (addr string, err error) {
	if len(prikey) != 32 {
		return
	}
	var s [32]byte
	copy(s[:], prikey)
	key, err := sr25519.NewMiniSecretKeyFromRaw(s)
	if err != nil {
		return "", err
	}
	pubK, err := key.ExpandEd25519().Public()
	if err != nil {
		return "", err
	}
	pub := pubK.Encode()
	addr, err = ss58.Encode(pub[:], prefix)
	if err != nil {
		return "", err
	}
	return addr, nil
}

// DotAddressFromPublicKey func
func DotAddressFromPublicKey(pubkey [32]byte, prefix []byte) (addr string, err error) {
	addr, err = ss58.Encode(pubkey[:], prefix)
	if err != nil {
		return "", err
	}
	return addr, nil
}

// DotAddressValidate func
func DotAddressValidate(addr string, prefix []byte) (err error) {
	err = ss58.VerityAddress(addr, prefix)
	return err
}

// DotPublicKey func
func DotPublicKey(prikey []byte) (pubkey string, err error) {
	if len(prikey) != 32 {
		return
	}
	var s [32]byte
	copy(s[:], prikey)
	key, err := sr25519.NewMiniSecretKeyFromRaw(s)
	if err != nil {
		return "", err
	}
	pubK, err := key.ExpandEd25519().Public()
	if err != nil {
		return "", err
	}
	pub := pubK.Encode()
	pubkey = hex.EncodeToString(pub[:])

	return pubkey, nil
}

// DotPublicKeyFromAddress func
func DotPublicKeyFromAddress(addr string) (pubkey string, err error) {
	pubkey = tx.AddressToPublicKey(addr)
	if pubkey == "" {
		return "", errors.New("bad address")
	}
	return pubkey, nil
}

// DotCreateSignedTransaction func
func DotCreateSignedTransaction(fromAddr, toAddr string,
	amount, nonce, fee uint64,
	specVer, txVer uint32, callIdx string,
	prik string,
	genesisHash, blockHash string, blockHeight uint64) (txid, txSigned string, err error) {

	btTx := tx.CreateTransaction(fromAddr, toAddr, amount, nonce, fee)
	btTx.SetGenesisHashAndBlockHash(genesisHash, blockHash, blockHeight)
	btTx.SetSpecVersionAndCallId(specVer, txVer, callIdx)
	_, txstr, err := btTx.CreateEmptyTransactionAndMessage()
	if err != nil {
		return "", "", err
	}
	sig, err := btTx.SignTransaction(prik, txstr)
	if err != nil {
		return "", "", err
	}
	txSigned, err = btTx.GetSignTransaction(sig)
	if err != nil {
		return "", "", err
	}
	txid, _, err = DotDecodeSignedTransaction(txSigned)
	if err != nil {
		return "", "", err
	}
	return txid, txSigned, nil
}

// DotDecodeSignedTransaction func
func DotDecodeSignedTransaction(txSigned string) (txid, txjson string, err error) {
	extrinsic := txSigned
	e := codes.ExtrinsicDecoder{}
	var option types.ScaleDecoderOption
	zbs, err := hex.DecodeString(config.OptionMetadata)
	if err != nil {
		return "", "", err
	}
	optbs, err := gu.BytesUnlzma(zbs)
	err = json.Unmarshal(optbs, &option)
	if err != nil {
		return "", "", err
	}
	e.Init(types.ScaleBytes{Data: iutiles.HexToBytes(extrinsic)}, &option)
	e.Process()
	txid = e.ExtrinsicHash
	bb, err := json.MarshalIndent(e.Value, "", "    ")
	if err != nil {
		return "", "", err
	}
	var resp v11.ExtrinsicDecodeResponse
	err = json.Unmarshal(bb, &resp)
	if err != nil {
		return "", "", err
	}
	txjson = string(bb)

	return txid, txjson, nil
}

// DotGetStorageFactor func
func DotGetStorageFactor(addr string) (fhex string, err error) {
	/*
		var hasher hash.Hash = xxhash.New128(nil)
		_, err = hasher.Write([]byte("Balances"))
		if err != nil {
			return "", err
		}
		bs := make([]byte, 0)
		bs = hasher.Sum(bs)
		ss := hex.EncodeToString(bs)
		fmt.Println(ss)
		// c2261276cc9d1f8598ea4b6a74b15c2f

		hasher.Reset()
		_, err = hasher.Write([]byte("FreeBalance"))
		if err != nil {
			return "", err
		}
		bs = make([]byte, 0)
		bs = hasher.Sum(bs)
		ss = hex.EncodeToString(bs)
		fmt.Println(ss)
		// 6482b9ade7bc6657aaca787ba1add3b4
	*/
	//fhex = "0xc2261276cc9d1f8598ea4b6a74b15c2f6482b9ade7bc6657aaca787ba1add3b4"
	fhex = "0xc2261276cc9d1f8598ea4b6a74b15c2f218f26c73add634897550b4003b26bc600"

	pubk, err := DotPublicKeyFromAddress(addr)
	if err != nil {
		return "", err
	}
	bs, err := hex.DecodeString(pubk)
	if err != nil {
		return "", err
	}

	hasher := blake2b.New256()
	_, err = hasher.Write(bs)
	if err != nil {
		return "", err
	}
	bs = make([]byte, 0)
	bs = hasher.Sum(bs)
	ss := hex.EncodeToString(bs)
	fhex = fhex + ss

	return fhex, nil
}
