package v1

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// NewGRPCHealthServer creates the standard grpc health server.
// It's tiny and requires no custom logic here.
func NewGRPCHealthServer(logger *zap.Logger) *health.Server {
	h := health.NewServer()
	logger.Info("created grpc health server")
	return h
}

// RegisterGRPCHealthServer registers the health server on the given grpc.Server.
func RegisterGRPCHealthServer(s *grpc.Server, h *health.Server) {
	healthpb.RegisterHealthServer(s, h)
}
