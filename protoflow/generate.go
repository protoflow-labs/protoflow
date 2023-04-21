package protoflow

//go:generate protoc --go_out=./ --twirp_out=./ -I./proto "./proto/workflow.proto"
