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

func author(uc *service.UserService) middleware.Middleware {
	keyPrivoder := func(token *jwt2.Token) (interface{}, error) {
		return jwt2.ParseRSAPublicKeyFromPEM([]byte(uc.GetJWTPK()))
	}
	return jwt.Server(
		keyPrivoder,
		jwt.WithSigningMethod(jwt2.SigningMethodRS256),
		jwt.WithClaims(func() jwt2.Claims { return &common.MyClaims{} }),
	)
}

func injectJwtInfo() middleware.Middleware {
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

var loginPath = []string{
	"/api.user.v1.User/LoginUser",
}

var withoutAuthPath = []string{
	"/api.user.v1.User/GetLoginInfo",
}

func needCheckJwt(ctx context.Context, operation string) bool {
	for _, v := range loginPath {
		if v == operation {
			return false
		}
	}
	for _, v := range withoutAuthPath {
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
			selector.Server(author(uc), injectJwtInfo()).Match(needCheckJwt).Build(),
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
