package xrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const scheme = "test"

type Config struct {
	Domain          string
	Endpoints       []string
	BackoffInterval time.Duration //设置默认回退策略的重试周期
	MaxAttempts     int           //默认为3，最大重连次数，若为负数则无限重连
}

// myResolver 自定义name resolver，实现Resolver接口
type myResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *myResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (*myResolver) Close() {}

// myResolverBuilder 需实现 Builder 接口
type myResolverBuilder struct {
	domain    string
	endpoints []string
}

func (m *myResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &myResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			m.domain: m.endpoints,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}
func (*myResolverBuilder) Scheme() string { return scheme }

// 拦截器实现重连
func retryInterceptor(maxAttempt int, backoffInterval time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		for attempt := 0; attempt <= maxAttempt || maxAttempt < 0; attempt++ {

			if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					log.Printf("Connection failed, retry attempt %d\n", attempt+1)
					log.Printf("err:%v", err)
					interval := time.Second
					if backoffInterval != 0 {
						interval = backoffInterval
					}
					time.Sleep(interval * time.Second) // 可以根据需要调整重试间隔
					continue
				}
				return err
			}
			return nil // 请求成功，不需要重试
		}
		return fmt.Errorf("Max retry attempts reached")
	}
}

func MustInitClient(config Config, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}

	if config.BackoffInterval == 0 {
		config.BackoffInterval = 1
	}

	opts = append(
		opts,
		grpc.WithUnaryInterceptor(retryInterceptor(config.MaxAttempts, config.BackoffInterval)),
		grpc.WithResolvers(&myResolverBuilder{config.Domain, config.Endpoints}),
	)

	conn, err := grpc.Dial(scheme+":///"+config.Domain, opts...)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	return conn, nil
}
