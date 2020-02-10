// Package proto contains generated protocol buffers.
package proto

//go:generate protoc -I=../../protocol --go_out=plugins=grpc:. login.proto
//go:generate protoc -I=../../protocol --go_out=plugins=grpc:. joyity.service.proto
