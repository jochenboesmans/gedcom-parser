#!/bin/bash
# compiles proto specs to go files
protoc --go_out=../../../. gedcom/gedcom.proto
protoc --go_out=../../../. grpc/parse.proto
