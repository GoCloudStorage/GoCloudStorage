package etcd

import (
	"context"
	"crypto/tls"

	"github.com/sirupsen/logrus"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	gResolver "google.golang.org/grpc/resolver"
)

var (
	Client   *etcdv3.Client
	Resolver gResolver.Builder
)

type Config struct {
	Context        context.Context
	Endpoints      []string
	TLS            *tls.Config
	Username       string `json:"username"`
	Password       string `json:"password"`
	EnableResolver bool
}

func MustNewClient(ctx context.Context, config Config) (*etcdv3.Client, gResolver.Builder) {
	var (
		err error
	)

	if Client, err = etcdv3.New(etcdv3.Config{
		Endpoints: config.Endpoints,
		TLS:       config.TLS,
		Username:  config.Username,
		Password:  config.Password,
		Context:   ctx,
	}); err != nil {
		logrus.Panicf("new etcd client err: %v", err)
	}

	if config.EnableResolver {
		if Resolver, err = resolver.NewBuilder(Client); err != nil {
			logrus.Panicf("new etcd resolver err: %v", err)
		}
	}

	logrus.Info("connect to etcd success")

	return Client, Resolver
}
