package data

import (
	"context"
	"fmt"
	userpb "realworld/api/user/v1"
	"realworld/app/mgr/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	srcgrpc "google.golang.org/grpc"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserProxy)

// Data .
type Data struct {
	// TODO wrapped database client

	userCli userpb.UserClient
}

// NewData .
func NewData(c *conf.Bootstrap, logger log.Logger) (*Data, func(), error) {
	d := Data{}
	connGRPC := initUserCli(c)
	d.userCli = userpb.NewUserClient(connGRPC)

	cleanup := func() {
		connGRPC.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &d, cleanup, nil
}

func initUserCli(c *conf.Bootstrap) (gc *srcgrpc.ClientConn) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: c.Data.Etcd.Addr,
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)
	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", c.UserSvc.SvcName)),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}

	return connGRPC
}
