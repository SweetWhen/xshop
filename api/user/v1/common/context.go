package common

import "context"

type jwtCtxKey struct{}

func SetJWTClaim(ctx context.Context, mc *MyClaims) context.Context {
	return context.WithValue(ctx, jwtCtxKey{}, mc)
}

func GetJWTClaim(ctx context.Context) (*MyClaims, bool) {
	v, ok := ctx.Value(jwtCtxKey{}).(*MyClaims)
	return v, ok
}
