# Metrics (Prometheus)

A wrapper for enabling [Prometheus](https://prometheus.io) on a stanard [http.ServerMux](https://golang.org/pkg/net/http/#ServeMux).

Also contains a docker-compose file for running a Prometheus server with a sample configuration.

## Usage

```go
mux := http.NewServeMux()
mux = metrics.Expose(mux)
```
