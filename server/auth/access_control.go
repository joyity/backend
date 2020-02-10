package auth

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type AccessControl struct {
	log *logrus.Logger
}

func NewAccessControl(log *logrus.Logger) *AccessControl {
	return &AccessControl{
		log: log,
	}
}

func (ac *AccessControl) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

func (ac *AccessControl) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(ctx, req)
	}
}
