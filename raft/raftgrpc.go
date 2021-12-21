package main

import (
	"context"
	"hash/fnv"
	"net"
	"strings"
	"sync"

	"example.com/m/pb"
	"github.com/gogo/protobuf/proto"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"google.golang.org/grpc"
)

type RaftGrpc struct {
	sync.RWMutex
	raftNode raft.Node
	addr     string
	peers    map[uint64]*peer
	errorCh  chan error
}

type raftGrpcServer struct {
	raftNode raft.Node
}

type peer struct {
	id     uint64
	addr   string
	conn   *grpc.ClientConn
	client pb.RaftClient
}

func newRaftGrpc(n *Node) *RaftGrpc {
	addr := n.peerMap[n.id][strings.LastIndex(n.peerMap[n.id], ":"):]
	return &RaftGrpc{addr: addr, raftNode: n.raftNode, errorCh: make(chan error)}
}

func (raftGrpc *RaftGrpc) Start() {
	listen, err := net.Listen("tcp", raftGrpc.addr)
	if err != nil {
		return
	}
	s := grpc.NewServer()
	pb.RegisterRaftServer(s, &raftGrpcServer{raftNode: raftGrpc.raftNode})
	go func() {
		err := s.Serve(listen)
		if err != nil {
			raftGrpc.errorCh <- err
		}
	}()
}

func (raftGrpc *RaftGrpc) AddPeer(node *pb.NodeInfo) error {
	conn, err := grpc.Dial(node.Addr)
	if err != nil {
		return err
	}
	client := pb.NewRaftClient(conn)
	p := peer{
		id:     node.Id,
		addr:   node.Addr,
		conn:   conn,
		client: client,
	}
	raftGrpc.Lock()
	defer raftGrpc.Unlock()
	raftGrpc.peers[p.id] = &p
	return nil
}

func (raftGrpc *RaftGrpc) Send(msgs []raftpb.Message) {
	peers := raftGrpc.getPeers()
	for _, m := range msgs {
		if p, ok := peers[m.To]; ok {
			req := pb.SendReq{
				Msg: &m,
			}
			_, err := p.client.Send(context.Background(), &req)
			if err != nil {
				raftGrpc.raftNode.ReportUnreachable(p.id)
			}
		}
	}
}

func (gt *RaftGrpc) getPeer(id uint64) *peer {
	gt.RLock()
	defer gt.RUnlock()
	return gt.peers[id]
}

func (gt *RaftGrpc) getPeers() map[uint64]*peer {
	gt.RLock()
	defer gt.RUnlock()
	ps := make(map[uint64]*peer, len(gt.peers))
	for k, v := range gt.peers {
		ps[k] = v
	}
	return ps
}

func (server *raftGrpcServer) Join(ctx context.Context, info *pb.NodeInfo) (resp *pb.Resp, err error) {
	id := genId(info.Addr)
	byt, _ := proto.Marshal(info)
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  id,
		Context: byt,
	}
	err = server.raftNode.ProposeConfChange(context.Background(), cc)
	if err != nil {
		return
	}
	resp = &pb.Resp{
		Success: true,
	}
	return
}

func (server *raftGrpcServer) Leave(ctx context.Context, info *pb.NodeInfo) (resp *pb.Resp, err error) {
	id := genId(info.Addr)
	byt, _ := proto.Marshal(info)
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeRemoveNode,
		NodeID:  id,
		Context: byt,
	}
	err = server.raftNode.ProposeConfChange(context.Background(), cc)
	if err != nil {
		return
	}
	resp = &pb.Resp{
		Success: true,
	}
	return
}

func (server *raftGrpcServer) Send(ctx context.Context, in *pb.SendReq) (resp *pb.Resp, err error) {
	err = server.raftNode.Step(ctx, *in.Msg)
	if err != nil {
		return
	}
	resp = &pb.Resp{
		Success: true,
	}
	return
}

func (raftGrpc *RaftGrpc) ErrorC() chan error {
	return raftGrpc.errorCh
}

func genId(addr string) uint64 {
	h := fnv.New64()
	h.Write([]byte(addr))
	return h.Sum64()
}
