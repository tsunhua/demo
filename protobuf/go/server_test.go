package main

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"testing"

	"example.com/m/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestWhoami(t *testing.T) {
	request, err := http.NewRequest("GET", "http://localhost:8080/whoami", nil)
	if err != nil {
		t.Error(err)
	}
	request.Header.Add("Content-Type", "application/x-protobuf")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 200, response.StatusCode)
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "080c1205446f6b6b79", hex.EncodeToString(bs))
	whoami := &pb.Animal{}
	proto.Unmarshal(bs, whoami)
	assert.Equal(t, int64(12), whoami.Id)
	assert.Equal(t, "Dokky", whoami.Name)
}
