package biz

import (
	"crypto/sha512"
	"fmt"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

type PWEncode struct {
	opt *password.Options
}

func NewPWEncode() *PWEncode {
	return &PWEncode{
		opt: &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New},
	}
}

func (pw *PWEncode) Encode(org string) string {
	salt, encodedPwd := password.Encode(org, pw.opt)
	return fmt.Sprintf("%s$%s", salt, encodedPwd)
}

func (pw *PWEncode) Decode(org, encoded string) error {
	passwordInfo := strings.Split(encoded, "$")
	if len(passwordInfo) != 2 {
		return fmt.Errorf("len(passwordInfo):%d", len(passwordInfo))
	}
	check := password.Verify(org, passwordInfo[0], passwordInfo[1], pw.opt)
	if !check {
		return fmt.Errorf("verify failed")
	}
	return nil
}
