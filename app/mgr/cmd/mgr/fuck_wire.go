package main

import (
	"realworld/app/mgr/internal/biz"
	"realworld/app/mgr/internal/conf"
	"realworld/app/mgr/internal/data"
	"realworld/app/mgr/internal/server"
	"realworld/app/mgr/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	bu := data.NewUserProxy(dataData)

	mgrUserUsecase := biz.NewMgrUserUsecase(bu, logger)
	mgrUserSvc := service.NewMgrUserSvc(mgrUserUsecase)
	httpServer := server.NewHTTPServer(confServer, mgrUserSvc, logger)
	app := newApp(logger, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
