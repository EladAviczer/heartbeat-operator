package prober

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GrpcProber struct {
	Target  string
	Timeout time.Duration
}

func NewGrpcProber(target string, timeout time.Duration) *GrpcProber {
	return &GrpcProber{Target: target, Timeout: timeout}
}

func (p *GrpcProber) Check() bool {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()

	conn, err := grpc.NewClient(p.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[gRPC] Client creation failed for %s: %v", p.Target, err)
		return false
	}
	defer conn.Close()

	client := healthpb.NewHealthClient(conn)
	resp, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Printf("[gRPC] HealthCheck failed on %s: %v", p.Target, err)
		return false
	}

	return resp.Status == healthpb.HealthCheckResponse_SERVING
}
