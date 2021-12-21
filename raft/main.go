package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"
	"unicode/utf8"

	"go.etcd.io/etcd/raft/v3"
)

func main() {
	id := flag.Uint64("id", 1, "node id")
	flag.Parse()
	log.Printf("node %v", *id)

	peerMap := map[uint64]string{
		1: "http://127.0.0.1:22210",
		2: "http://127.0.0.1:22211",
		3: "http://127.0.0.1:22212",
	}
	n := NewNode(*id, peerMap)
	n.Start()

	utf8.RuneCountInString("好你")

	for {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		if n.raftNode.Status().RaftState == raft.StateLeader {
			log.Printf("propose, node: %v", *id)
			n.raftNode.Propose(context.TODO(), []byte("hello"))
		}
	}
}
