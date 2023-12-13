package server

import (
	"context"
	userpb "realworld/api/user/v1"
	"realworld/api/user/v1/common"
	"realworld/app/user/internal/conf"
	"realworld/app/user/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"
)

func loginCheck(uc *service.UserService) middleware.Middleware {
	keyProvider := func(t *jwt2.Token) (interface{}, error) { return []byte(uc.GetJWTSK()), nil }
	emptyHandler := func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil }
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			resp, err := h(ctx, req)
			if err != nil {
				return resp, err
			}
			loginResp := resp.(*userpb.LoginUserResp)
			//generate jwt
			jwtInjector := jwt.Client(keyProvider, jwt.WithSigningMethod(jwt2.SigningMethodRS256),
				jwt.WithClaims(func() jwt2.Claims { return &common.MyClaims{Uid: loginResp.Uid, Name: loginResp.Name} }))
			_, e := jwtInjector(emptyHandler)(ctx, req)
			if e != nil {
				return nil, e
			}
			return resp, nil
		}
	}
}

func author(uc *service.UserService) middleware.Middleware {
	keyPrivoder := func(token *jwt2.Token) (interface{}, error) {
		return []byte(uc.GetJWTPK()), nil
	}
	return jwt.Server(
		keyPrivoder,
		jwt.WithSigningMethod(jwt2.SigningMethodRS256),
		jwt.WithClaims(func() jwt2.Claims { return &common.MyClaims{} }),
	)
}

var loginPath = []string{
	"/api.user.v1.User/LoginUser",
}

func needCheckJwt(ctx context.Context, operation string) bool {
	for _, v := range loginPath {
		if v == operation {
			return false
		}
	}
	return true
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, uc *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server(loginCheck(uc)).Path(loginPath...).Build(),
			selector.Server(author(uc)).Match(needCheckJwt).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	userpb.RegisterUserHTTPServer(srv, uc)
	return srv
}
