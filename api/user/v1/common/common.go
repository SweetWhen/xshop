package common

import jwt2 "github.com/golang-jwt/jwt/v4"

type MyClaims struct {
	jwt2.RegisteredClaims
	Uid  int64
	Name string
}
