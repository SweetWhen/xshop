package biz

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"hash"
	"io"
)

type RSAImpl struct {
	privateObj            *rsa.PrivateKey
	publicKey, privateKey string
	h                     crypto.Hash
	encryHash             hash.Hash
	encryRand             io.Reader
}

func NewRSAImpl(sk string) (impl *RSAImpl) {
	impl = &RSAImpl{
		privateKey: sk,
		h:          crypto.SHA256,
		encryHash:  sha256.New(),
		encryRand:  rand.Reader,
	}
	if len(impl.privateKey) == 0 {
		if err := impl.initKeys(); err != nil {
			panic(fmt.Sprintf("initKeys err:%s", err.Error()))
		}
	} else {
		if err := impl.initObj(); err != nil {
			panic(fmt.Sprintf("initObj err:%s", err.Error()))
		}
	}

	return impl
}

func (impl *RSAImpl) initKeys() (err error) {
	impl.privateObj, err = rsa.GenerateKey(impl.encryRand, 1024)
	if err != nil {
		return
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(impl.privateObj),
	}
	impl.privateKey = string(pem.EncodeToMemory(privateKeyPEM))

	publicKey := &impl.privateObj.PublicKey
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	}
	impl.publicKey = string(pem.EncodeToMemory(publicKeyPEM))
	return
}

func (impl *RSAImpl) initObj() (err error) {
	pemBlock, _ := pem.Decode([]byte(impl.privateKey))
	if impl.privateObj, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err != nil {
		return
	}

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&impl.privateObj.PublicKey),
	}
	publicKeyBuff := &bytes.Buffer{}
	pem.Encode(publicKeyBuff, publicKeyPEM)
	impl.publicKey = publicKeyBuff.String()
	return
}

func (impl *RSAImpl) EncryAndBase64(msg string) (encryMsg string, err error) {
	publicKey := &impl.privateObj.PublicKey
	var encryMsgBytes, eachEncryBytes []byte
	h := impl.encryHash
	r := impl.encryRand
	orgMsg := []byte(msg)
	step := publicKey.Size() - 2*h.Size() - 2
	for start := 0; start < len(orgMsg); start += step {
		finish := start + step
		if finish > len(orgMsg) {
			finish = len(orgMsg)
		}
		eachEncryBytes, err = rsa.EncryptOAEP(h, r, publicKey, orgMsg[start:finish], nil)
		if err != nil {
			return
		}
		encryMsgBytes = append(encryMsgBytes, eachEncryBytes...)
	}

	encryMsg = base64.StdEncoding.EncodeToString(encryMsgBytes)
	return
}

func (impl *RSAImpl) DecryBase64(encry string) (decryMsg string, err error) {
	privateKey := impl.privateObj
	var decodeMsg []byte
	decodeMsg, err = base64.StdEncoding.DecodeString(encry)
	if err != nil {
		err = fmt.Errorf("base64 decode err:%s", err.Error())
		return
	}
	step := privateKey.PublicKey.Size()
	var decryBytes, decryEachBytes []byte

	for start := 0; start < len(decodeMsg); start += step {
		finish := start + step
		if finish > len(decodeMsg) {
			finish = len(decodeMsg)
		}
		decryEachBytes, err = privateKey.Decrypt(nil, decodeMsg[start:finish], &rsa.OAEPOptions{Hash: impl.h})
		if err != nil {
			err = fmt.Errorf("decrypt err:%s", err.Error())
			return
		}
		decryBytes = append(decryBytes, decryEachBytes...)
	}

	decryMsg = string(decryBytes)
	return
}
