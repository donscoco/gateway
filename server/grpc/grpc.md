
安装
参考官方：https://github.com/grpc/grpc-go
```bash
go mod export GO111MODULE=on
go mod export GOPROXY=https://goproxy.io
## 安装grpc
go get -u google.glang.org/grpc
## 安装proto
go get -u github.com/golang/protobuf/proto
## 安装protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

```