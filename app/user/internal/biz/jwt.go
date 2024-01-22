package biz

import (
	"context"
	"fmt"
	"time"

	userpb "realworld/api/user/v1"
	"realworld/api/user/v1/common"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"google.golang.org/protobuf/proto"
)

var jwtMap = map[string][]byte{}

func (uc *UserUsecase) Checkin(ctx context.Context, p *userpb.ClaimPayload) (jwtID string, err error) {
	//todo: 对称加密： account + sid过期时间戳
	// 另外用redis hset保存某个用户所有的sid
	jwtID = fmt.Sprintf("uid-%d", p.Uid)
	savePayload, e := proto.Marshal(p)
	if e != nil {
		err = e
		return
	}

	jwtMap[jwtID] = savePayload

	return
}

func (uc *UserUsecase) logout(ctx context.Context, uid int64) error {
	jwtID := fmt.Sprintf("uid-%d", uid)
	delete(jwtMap, jwtID)
	return nil
}

func (uc *UserUsecase) InjectJwt(claim *common.MyClaims) (header, value string, e error) {
	token := jwt2.NewWithClaims(jwt2.SigningMethodRS256, claim)
	tokenStr, err := token.SignedString(uc.rsaImpl.privateObj)
	if err != nil {
		e = err
		return
	}
	header = "Authorization"
	value = fmt.Sprintf("Bearer %s", tokenStr)
	return
}

func (uc *UserUsecase) SidCheck(sid string) (header, value string, err error) {
	jwtPayload, ok := jwtMap[sid]
	if !ok {
		err = fmt.Errorf("sid: %s not found", sid)
		return
	}
	claim := &common.MyClaims{
		RegisteredClaims: jwt2.RegisteredClaims{
			ID:        sid,
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(uc.confServer.JwtInterval.AsDuration())),
		},
		ClaimPayload: userpb.ClaimPayload{},
	}
	err = proto.Unmarshal(jwtPayload, &claim.ClaimPayload)
	if err != nil {
		return
	}
	header, value, err = uc.InjectJwt(claim)
	return
}

func (uc *UserUsecase) JwtCheck(tokenStr string) (*common.MyClaims, error) {
	keyPrivoder := func(token *jwt2.Token) (interface{}, error) {
		return &uc.rsaImpl.privateObj.PublicKey, nil
	}
	reason := "UNAUTHORIZED"
	claims := &common.MyClaims{}
	token, err := jwt2.ParseWithClaims(tokenStr, claims, keyPrivoder)
	if err != nil {
		ve, ok := err.(*jwt2.ValidationError)
		if !ok {
			return nil, errors.Unauthorized(reason, err.Error())
		}
		if ve.Errors&jwt2.ValidationErrorMalformed != 0 {
			return nil, fmt.Errorf("ErrTokenInvalid")
		}
		if ve.Errors&(jwt2.ValidationErrorExpired|jwt2.ValidationErrorNotValidYet) != 0 {
			return nil, fmt.Errorf("ErrTokenExpired")
		}
		if ve.Inner != nil {
			fmt.Printf("JwtCheck inner err:%v\n", ve.Inner)
			return nil, ve.Inner
		}
		return nil, fmt.Errorf("ErrTokenParseFail")
	}
	if !token.Valid {
		return nil, fmt.Errorf("ErrTokenInvalid")
	}

	return claims, nil

}

func (uc *UserUsecase) ParserJwtInfo() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if t, ok := jwt.FromContext(ctx); ok {
				c := t.(*common.MyClaims)
				ctx = common.SetJWTClaim(ctx, c)
			}
			return h(ctx, req)
		}
	}
}
