package xrpc

import (
	"google.golang.org/grpc"
	"sync"
)

var (
	mu      sync.RWMutex
	clients = make(map[string]interface{})
)

type MidClient[T any] struct {
	Config
	opts []grpc.DialOption
	cc   *grpc.ClientConn
	fn   func(connInterface grpc.ClientConnInterface) T
}

func (c *MidClient[T]) NewSession() T {
	return c.fn(c.cc)
}

func (c *MidClient[T]) dail() error {
	var err error
	c.cc, err = MustInitClient(c.Config, c.opts...)
	return err
}

func GetGrpcClient[T any](config Config, fn func(grpc.ClientConnInterface) T, opts ...grpc.DialOption) (*MidClient[T], error) {
	mu.RLock()
	v, exist := clients[config.Domain]
	mu.RUnlock()
	if exist && v.(*MidClient[T]).cc != nil {
		return v.(*MidClient[T]), nil
	}

	mu.Lock()
	defer mu.Unlock()
	client := &MidClient[T]{
		Config: config,
		fn:     fn,
		opts:   opts,
	}
	err := client.dail()
	if err != nil {
		return nil, err
	}
	clients[config.Domain] = clients

	return client, nil
}
