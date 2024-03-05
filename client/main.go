package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/muzzapp/devops-interview-task/pkg/muzz"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_calls_total",
		Help: "Total number of gRPC calls",
	}, []string{"method", "status"})

	grpcLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "grpc_call_latency_seconds",
		Help:    "Latency of gRPC calls",
		Buckets: prometheus.DefBuckets,
	})
)

func main() {
	conn, dialErr := grpc.Dial("127.0.0.1:9876", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if dialErr != nil {
		slog.Error("Failed to dial gRPC service: %v", "err", dialErr)
		os.Exit(1)
	}

	client := muzz.NewServiceClient(conn)

	start := time.Now()
	response, respErr := client.Echo(
		context.Background(),
		&muzz.EchoRequest{Message: "Hello, world!"},
	)
	elapsed := time.Since(start)
	grpcLatency.Observe(elapsed.Seconds())

	if respErr != nil {
		slog.Error("Failed to call gRPC service: %v", "err", respErr)
		grpcCalls.With(prometheus.Labels{"method": "Echo", "status": "fail"}).Inc()
		os.Exit(1)
	}

	slog.Info("Response: %s", "message", response.Message)
	grpcCalls.With(prometheus.Labels{"method": "Echo", "status": "success"}).Inc()
}
