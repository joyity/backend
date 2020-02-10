package server

import "github.com/joyity/backend/server/proto"

var _ proto.JoyityServer = (*protoService)(nil)

type protoService struct {
	proto.UnimplementedJoyityServer
}

func newProtoService() *protoService {
	return &protoService{}
}
