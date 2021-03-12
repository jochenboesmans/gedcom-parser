#!/bin/bash
# compiles proto specs to go files
# requires protoc to be installed
go get google.golang.org/protobuf/cmd/protoc-gen-go google.golang.org/grpc/cmd/protoc-gen-go-grpc
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=../../../. --go-grpc_out=../../../. gedcom/gedcom.proto
protoc --go_out=../../../. --go-grpc_out=../../../. grpc/parse.proto
