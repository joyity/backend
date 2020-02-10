package server

import (
	"fmt"
	"io"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/joyity/backend/server/auth"
	"github.com/joyity/backend/server/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server interface {
	ListenAndServe() error

	io.Closer
}

func New(log *logrus.Logger) Server {
	return newServer(log)
}

// server implementation

var _ Server = (*server)(nil)

type server struct {
	log *logrus.Logger

	accessControl *auth.AccessControl
	grpcServer    *grpc.Server
	protoService  proto.JoyityServer
}

func newServer(log *logrus.Logger) *server {
	accessControl := auth.NewAccessControl(log.WithField("component", "access control").Logger)

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			accessControl.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(log).WithField("component", "grpc")),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			accessControl.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(log).WithField("component", "grpc")),
		)),
	)
	protoService := newProtoService()

	proto.RegisterJoyityServer(grpcServer, protoService)

	return &server{
		log: log,

		accessControl: accessControl,
		grpcServer:    grpcServer,
		protoService:  protoService,
	}
}

func (s *server) ListenAndServe() error {
	protocol := "tcp"
	addr := ":50432"
	lis, err := net.Listen(protocol, addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	s.log.Infof("listening on %v (%v)", addr, protocol)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

func (s *server) Close() error {
	// Try to gracefully stop the server. If that doesn't work within 3 seconds,
	// stop the server more forcefully.
	timeout := 3 * time.Second
	gracefullyStoppedSignal := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(gracefullyStoppedSignal)
	}()
	select {
	case <-gracefullyStoppedSignal:
	case <-time.After(timeout):
		s.grpcServer.Stop()
	}

	return nil
}
