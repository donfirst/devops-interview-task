package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/muzzapp/devops-interview-task/pkg/muzz"
	"github.com/muzzapp/devops-interview-task/server/internal/handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	totalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests received",
		},
	)
	successfulRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "successful_requests",
			Help: "Total number of successful requests",
		},
	)
	failedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "failed_requests",
			Help: "Total number of failed requests",
		},
	)
)

func init() {
	prometheus.MustRegister(totalRequests, successfulRequests, failedRequests)
}
func main() {
	// Create a new gRPC server
	server := grpc.NewServer()
	muzz.RegisterServiceServer(server, handler.Server{})

	// Run the gRPC server
	go func() {
		listener, listenerErr := net.Listen("tcp", "0.0.0.0:9876")
		if listenerErr != nil {
			slog.Error("Failed to create gRPC listener", "err", listenerErr)
			os.Exit(1)
		}

		slog.Info("gRPC server starting")
		if err := server.Serve(listener); err != nil {
			slog.Error("gRPC server stopped", "err", err)
			os.Exit(1)
		}
	}()

	// Run the Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	// Gracefully shut down gRPC server after receiving an interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	slog.Info("Shutting down gRPC server")
	server.GracefulStop()
}
