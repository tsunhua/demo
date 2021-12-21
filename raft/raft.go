package main

import (
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"log"
	"time"
)

type Node struct {
	id            uint64
	peerMap       map[uint64]string
	raftNode      raft.Node
	raftStorage   *raft.MemoryStorage
	raftTransport RaftTransport
}

type RaftTransport interface {
	Start()
	Send(msgs []raftpb.Message)
	ErrorC() chan error
}

func NewNode(id uint64, peerMap map[uint64]string) *Node {
	n := &Node{
		id:          id,
		peerMap:     peerMap,
		raftStorage: raft.NewMemoryStorage(),
	}
	return n
}

func (n *Node) Start() {
	peers := make([]raft.Peer, 0, len(n.peerMap))
	for i := range n.peerMap {
		peers = append(peers, raft.Peer{ID: uint64(i)})
	}
	c := &raft.Config{
		ID:              n.id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         n.raftStorage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	n.raftNode = raft.StartNode(c, peers)
	// n.raftServer = newRaftHttp(n)
	n.raftTransport = newRaftGrpc(n)
	n.raftTransport.Start()
	go n.doBackgroundWork()
}

func (n *Node) doBackgroundWork() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.raftNode.Tick()
		case rd := <-n.raftNode.Ready():
			n.raftStorage.Append(rd.Entries)
			n.raftTransport.Send(rd.Messages)
			if !raft.IsEmptySnap(rd.Snapshot) { // snapshot 何处来？
				n.raftStorage.ApplySnapshot(rd.Snapshot)
			}
			for _, entry := range rd.CommittedEntries {
				switch entry.Type {
				case raftpb.EntryNormal:
				case raftpb.EntryConfChange:
					var cc raftpb.ConfChange
					cc.Unmarshal(entry.Data)
					n.raftNode.ApplyConfChange(cc)
				}
			}
			n.raftNode.Advance()
		case err := <-n.raftTransport.ErrorC():
			log.Fatal(err)
		}
	}
}
