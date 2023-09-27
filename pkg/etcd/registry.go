package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"net/url"
	"os"
	"time"
)

type Service struct {
	Project      string    `json:"project"`       // default: bifrost-pro-local
	Name         string    `json:"name"`          // topic, content, search ...
	MsgType      string    `json:"msg_type"`      // registry, unregistry
	Scheme       string    `json:"scheme"`        // http, https, default: http
	Host         string    `json:"address"`       // 10.220.10.15, default: name
	Port         int       `json:"port"`          // 34001, default: 80
	Replica      string    `json:"replica"`       // default: env(HOSTNAME)
	RegTime      time.Time `json:"reg_time"`      // default: time.Now()
	AuthExcludes []string  `json:"auth_excludes"` // /{{module_name}}/api_prefix_1/api_prefix_2/api
}

func (s *Service) ToString() string {
	if bs, err := json.Marshal(s); err != nil {
		logrus.Panicf("marshal service to bytes err: %v", err)
	} else {
		return string(bs)
	}

	return ""
}

func MustRegistry(gctx context.Context, etcdClient *etcdv3.Client, srv Service) {
	var err error
	if srv.Host == "" {
		logrus.Warnf("service host invalid, set to name: %v", srv.Name)
		srv.Host = srv.Name
	}

	if srv.Project == "" {
		logrus.Warn("service project empty, set default(bifrost-pro-local)")
		srv.Project = "bifrost-pro-local"
	}

	if srv.Replica == "" {
		logrus.Warnf("service replica invalid, set default(HOSTNAME: %s)", os.Getenv("HOSTNAME"))
		srv.Replica = os.Getenv("HOSTNAME")
	}

	if srv.RegTime.IsZero() {
		now := time.Now()
		logrus.Warnf("service replica invalid, set default(%v)", now)
		srv.RegTime = now
	}

	if srv.MsgType == "" {
		srv.MsgType = "registry"
	}

	if srv.Scheme == "" {
		srv.Scheme = "http"
	}

	// 开始验证
	if srv.Host == "" || srv.Name == "" {
		logrus.Panic("must declare service address,name")
	}

	if srv.MsgType != "registry" && srv.MsgType != "unregistry" {
		logrus.Panic("service msg_Type only accept (registry, unregistry)")
	}

	if _, err = url.Parse(fmt.Sprintf("%s://%s:%d", srv.Scheme, srv.Host, srv.Port)); err != nil {
		logrus.Panicf("service address invalid: %v", err)
	}

	var (
		grantResp             *etcdv3.LeaseGrantResponse
		leaseCh               <-chan *etcdv3.LeaseKeepAliveResponse
		manager               endpoints.Manager
		hostname              = fmt.Sprintf("resolver/%s/%s", srv.Project, srv.Name)
		replicaName           = fmt.Sprintf("resolver/%s/%s/%s", srv.Project, srv.Name, srv.Replica)
		leaseCtx, leaseCancel = context.WithCancel(context.Background())
	)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if grantResp, err = etcdClient.Grant(ctx, 20); err != nil {
		logrus.Panicf("etcd client grant err: %v", err)
	}
	lease := etcdv3.NewLease(etcdClient)
	if leaseCh, err = lease.KeepAlive(leaseCtx, grantResp.ID); err != nil {
		logrus.Panicf("etcd client keep alive(get ka ch) err: %v", err)
	}

	if manager, err = endpoints.NewManager(etcdClient, hostname); err != nil {
		logrus.Panicf("etcd client new manager err: %v", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	if err = manager.AddEndpoint(
		ctx,
		replicaName,
		endpoints.Endpoint{
			Addr:     fmt.Sprintf("%s:%d", srv.Host, srv.Port),
			Metadata: srv.ToString(),
		},
		etcdv3.WithLease(grantResp.ID),
	); err != nil {
		logrus.Panicf("etcd manager add endpoint err: %v", err)
	}

	go func() {
		for {
			select {
			case msg, ok := <-leaseCh:
				if !ok {
					logrus.Panicf("etcd keep alive lease chan closed, service: %s", srv.ToString())
					return
				}

				_ = msg
			case <-gctx.Done():
				logrus.Warnf("etcd keep alive lease chan closed by gctx, service: %s", srv.ToString())
				Unregistry(etcdClient, srv)
				leaseCancel()
				return
			}
		}
	}()
	logrus.Infof("etcd registry(%s/%s/%s) success! detail: %s", srv.Project, srv.Name, srv.Replica, srv.ToString())

}

func Unregistry(client *etcdv3.Client, srv Service) {
	if srv.Project == "" {
		logrus.Warn("service project empty, set default(bifrost-pro-local)")
		srv.Project = "bifrost-pro-local"
	}

	if srv.Replica == "" {
		logrus.Warn("service replica invalid, set default(1)")
		srv.Replica = os.Getenv("HOSTNAME")
	}

	srv.MsgType = "unregistry"

	var (
		err         error
		manager     endpoints.Manager
		hostname    = fmt.Sprintf("resolver/%s/%s", srv.Project, srv.Name)
		replicaName = fmt.Sprintf("resolver/%s/%s/%s", srv.Project, srv.Name, srv.Replica)
		lr          *etcdv3.LeaseLeasesResponse
	)

	timeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
	lr, err = client.Leases(timeout)
	if err != nil {
		logrus.Warnf("service: %s unregistry new lease err: %v", replicaName, err)
	} else if len(lr.Leases) == 0 {
		logrus.Warnf("service: %s unregistry new lease length == 0", replicaName)
	} else {
		bs, _ := json.Marshal(map[string]any{"Metadata": srv.ToString()})
		if _, err = client.Put(timeout, replicaName, string(bs), etcdv3.WithLease(lr.Leases[0].ID)); err != nil {
			logrus.Warnf("service: %s unregistry put err: %v", replicaName, err)
		}
	}

	if manager, err = endpoints.NewManager(client, hostname); err != nil {
		logrus.Errorf("etcd client new manager err: %v", err)
		return
	}

	if err = manager.DeleteEndpoint(timeout, replicaName); err != nil {
		logrus.Errorf("etcd client delete endpoint err: %v", err)
		return
	}
	logrus.Infof("etcd deleted endpoint: %s", srv.ToString())

}
