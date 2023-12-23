package common

import (
	userpb "realworld/api/user/v1"

	jwt2 "github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	jwt2.RegisteredClaims
	userpb.ClaimPayload
}
