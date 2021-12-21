# raft

## 生成 protobuf文件

```
protoc --proto_path=. \
      --proto_path=./vendor \
      --proto_path=./vendor/github.com/gogo/protobuf \
      --proto_path=$GOPATH/bin/protoc-3.15.8/include \
      --gofast_out=plugins=grpc:. \
      pb/raft_grpc.proto
```
