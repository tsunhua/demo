package main

import (
	"encoding/hex"

	"example.com/m/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func main() {
	s := gin.Default()
	s.GET("/whoami", func(c *gin.Context) {
		animal := pb.Animal{
			Id:   12,
			Name: "Dokky",
		}
		bs, _ := proto.Marshal(&animal)
		println(hex.EncodeToString(bs))
		c.Data(200, "application/x-protobuf", bs)
	})
	s.Run()
}
