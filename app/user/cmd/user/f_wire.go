package main

import (
	"realworld/app/user/internal/biz"
	"realworld/app/user/internal/conf"
	"realworld/app/user/internal/data"
	"realworld/app/user/internal/server"
	"realworld/app/user/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp2(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(confServer, userRepo, logger, confData)
	userService := service.NewUserService(logger, userUsecase)
	grpcServer := server.NewGRPCServer(confServer, userService, logger)
	httpSvr := server.NewHTTPServer(confServer, userService, logger)
	app := newApp(logger, confData, grpcServer, httpSvr)
	return app, func() {
		cleanup()
	}, nil
}
