package biz

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {
	impl := NewRSAImpl("")
	privateKeyStr := impl.privateKey
	publicKeyStr := impl.publicKey
	fmt.Printf("privateKey: %s\n", privateKeyStr)
	fmt.Printf("publicKey: %s\n", publicKeyStr)

	newImpl := NewRSAImpl(privateKeyStr)

	if privateKeyStr != newImpl.privateKey {
		t.Fatalf("privateKeyStr != privateKey1, privateKeyStr:\n%s\n privateKey1:\n%s\n", privateKeyStr, newImpl.privateKey)
	}

	if publicKeyStr != newImpl.publicKey {
		t.Fatalf("publicKeyStr != publicKey1, publicKeyStr:\n%s\npublicKey1:\n%s\n  ", publicKeyStr, newImpl.publicKey)
	}
	msg := strings.Repeat("hello abcd", 5)
	ts := []struct {
		Msg string
	}{
		{Msg: msg},
		{Msg: strings.Repeat(msg, 2)},
		{Msg: strings.Repeat(msg, 3)},
		{Msg: strings.Repeat(msg, 5)},
		{Msg: strings.Repeat(msg, 7)},
		{Msg: strings.Repeat(msg, 20)},
	}
	for _, m := range ts {
		encryMsg, err := newImpl.EncryAndBase64(m.Msg)
		if err != nil {
			t.Fatalf("EncryAndBase64 err:%s", err.Error())
		}
		decryMsg, err := newImpl.DecryBase64(encryMsg)
		if err != nil {
			t.Fatalf("DecryBase64 err:%s", err.Error())
		}

		if reflect.DeepEqual(encryMsg, decryMsg) {
			t.Fatalf("encryMsg != decryMsg,encryMsg: %s,  decryMsg:%s", string(encryMsg), string(decryMsg))
		}

		fmt.Printf("encryMsg: len:%d, %s, decryMsg:%s\n",
			len(string(encryMsg)), string(encryMsg), string(decryMsg))
	}

}
