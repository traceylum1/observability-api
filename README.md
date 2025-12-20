# Observable Go API

A production-ready Go HTTP API demonstrating how to **embed observability directly into an application server** using structured logging, metrics, distributed tracing, and graceful shutdown.

This project is intentionally light on business logic and heavy on **production concerns**, serving as a **reference implementation / template** for real-world backend microservices.

## Project Goals

The primary goal of this project is to demonstrate how a Go microservice should be **structured and instrumented for production**, including:

- End-to-end request observability  
- Clean middleware design  
- Context propagation  
- Safe error handling  
- Graceful shutdown  

Each request produces **correlated logs, metrics, and traces**, allowing operators to answer:

- **Is the service healthy?** (metrics)  
- **Where is time being spent?** (traces)  
- **Why did this request fail?** (logs)  

---

## API Endpoints

| Method | Path          | Description                      |
|------:|---------------|----------------------------------|
| GET   | `/healthz`    | Liveness check                   |
| GET   | `/readyz`     | Readiness check                  |
| GET   | `/items/{id}` | Simulated read workload          |
| POST  | `/items`      | Simulated write workload         |
| GET   | `/metrics`    | Prometheus metrics endpoint      |

---

## Observability Features

### Structured Logging
- JSON-formatted logs  
- Request-scoped fields:
  - request ID  
  - method  
  - path  
  - status code  
  - latency  
- Panic stack traces captured safely  

---

### Metrics (Prometheus)

Exposed metrics include:
- HTTP request count (by method, route, status)
- Request latency histograms
- In-flight request gauge
- Panic counter

Metrics are labeled using **route templates** to avoid cardinality explosion.

---

### Distributed Tracing (OpenTelemetry)
- Root span per incoming HTTP request
- Child spans for handler execution and simulated dependencies
- Context propagation via HTTP headers
- Ready to export to Jaeger, Tempo, or other OTLP-compatible backends

---

### Panic Recovery
- Recovers panics at the middleware layer
- Logs structured error details and stack traces
- Increments panic metrics
- Returns a safe `500 Internal Server Error`

---

### Graceful Shutdown
- Handles `SIGINT` and `SIGTERM`
- Stops accepting new requests
- Allows in-flight requests to complete
- Shuts down with a configurable timeout  

This behavior is critical for containerized environments such as **Kubernetes**, **Cloud Run**, and **ECS**.

