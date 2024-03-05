# Devops Interview Task

This repository is used as part of the interview process for a DevOps Engineer at Muzz.

You'll need the latest [Go runtime installed](https://go.dev/dl/) to get started.

To run the server run:

```bash
go run server/main.go
```

To run the client run:

```bash
go run client/main.go
```


----------------
# Server

The metrics are collected using Prometheus, open-source system monitoring and alerting toolkit. The metrics include
  total requests,
  successful requests,
  failed requests.
A new HTTP server is also added for Prometheus to scrape metrics from.

This code is sending metrics to Prometheus by using the Prometheus Go client library
# Metrics Definition

The code defines two metrics:

grpc_calls_total: This is a counter vector that tracks the total number of gRPC calls. It includes labels for the method name and the status of the call (success or failure).

grpc_call_latency_seconds: This is a histogram that measures the latency of gRPC calls in seconds.

Both metrics are created using the promauto.NewCounterVec and promauto.NewHistogram functions from the Prometheus client library.


# Metrics Collection

The code collects metrics within the main function:

Latency Measurement: It records the start time of the gRPC call using time.Now(). After the call, it calculates the elapsed time with time.Since(start). The elapsed time, in seconds, is then observed by the grpcLatency histogram using grpcLatency.Observe(elapsed.Seconds()).

RPC Call Count: After the gRPC call, the code checks if there was an error. If there was, it increments the grpc_calls_total counter for the "Echo" method with a "fail" status. If the call was successful, it increments the counter with a "success" status.


# Client

This code is sending metrics to Prometheus by using the Prometheus Go client library

The metrics added include a counter for successful and failed gRPC calls, and a histogram to measure the latency of the gRPC calls.

Here's how it works:

# Metrics Definition

Two metrics are defined at the beginning of the code:

grpc_calls_total: This is a counter vector that tracks the total number of gRPC calls. It has labels for the method name and the status of the call (success or failure). The promauto.NewCounterVec function is used to create this metric.

grpc_call_latency_seconds: This is a histogram that measures the latency of gRPC calls in seconds. The promauto.NewHistogram function is used to create this metric.

# Metrics Collection

The metrics are collected in the main function:

Latency Measurement: The start time of the gRPC call is recorded using time.Now(). After the call, the elapsed time is calculated with time.Since(start). This elapsed time, in seconds, is observed by the grpcLatency histogram with grpcLatency.Observe(elapsed.Seconds()).

gRPC Call Count: After the gRPC call, the code checks if there was an error. If there was, it increments the grpc_calls_total counter for the "Echo" method with a "fail" status using grpcCalls.With(prometheus.Labels{"method": "Echo", "status": "fail"}).Inc(). If the call was successful, it increments the counter with a "success" status using grpcCalls.With(prometheus.Labels{"method": "Echo", "status": "success"}).Inc().


# Metrics Exposure

The metrics are automatically exposed on a /metrics HTTP endpoint by the Prometheus Go client library. Prometheus can then scrape this endpoint to collect the metrics.

Note that the code for starting the HTTP server to expose the metrics is not included in this code intentionally.

The metrics won't be sent to Prometheus directly by this code. Instead, Prometheus will pull (or scrape) the metrics from the application. You need to configure Prometheus to know where to scrape the metrics from, which is typically done in the Prometheus configuration file.
