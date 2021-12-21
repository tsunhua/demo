package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/spf13/cast"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
)

type RaftHttp struct {
	addr      string
	transport *rafthttp.Transport
}

func newRaftHttp(n *Node) *RaftHttp {
	raftTransport := &rafthttp.Transport{
		ID:          types.ID(n.id),
		ClusterID:   0x1000,
		Raft:        &raftWrapper{n.raftNode},
		ServerStats: stats.NewServerStats("server"+cast.ToString(n.id), cast.ToString(n.id)),
		LeaderStats: stats.NewLeaderStats(&zap.Logger{}, cast.ToString(n.id)),
		ErrorC:      make(chan error),
	}
	raftTransport.Start()
	for id, addr := range n.peerMap {
		if id != n.id {
			raftTransport.AddPeer(types.ID(id), []string{addr})
		}
	}
	addr := n.peerMap[n.id][strings.LastIndex(n.peerMap[n.id], ":"):]
	return &RaftHttp{transport: raftTransport, addr: addr}
}

func (server *RaftHttp) Start() {
	go func() {
		server := http.Server{
			Addr:    server.addr,
			Handler: server.transport.Handler(),
		}
		server.ListenAndServe()
	}()
}

func (server *RaftHttp) Send(msgs []raftpb.Message) {
	server.transport.Send(msgs)
}

func (server *RaftHttp) ErrorC() chan error {
	return server.transport.ErrorC
}

type raftWrapper struct {
	node raft.Node
}

func (wrapper *raftWrapper) Process(ctx context.Context, m raftpb.Message) error {
	return wrapper.node.Step(ctx, m)
}

func (wrapper *raftWrapper) IsIDRemoved(id uint64) bool {
	return false
}

func (wrapper *raftWrapper) ReportUnreachable(id uint64) {
	wrapper.node.ReportUnreachable(id)
}

func (wrapper *raftWrapper) ReportSnapshot(id uint64, status raft.SnapshotStatus) {
	wrapper.node.ReportSnapshot(id, status)
}
