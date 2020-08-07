An example of generating type-safe Envoy configuration in Go using Envoy's protocol buffers.

### Build the configuration and run the envoy server in Docker:

```
docker build -f Dockerfile-envoy --tag envoy-config:latest .
```

### Run the envoy server:

```
docker run -d -p 10000:10000 envoy:latest
```

Confirm that all is working by issuing a request:

```
curl 0.0.0.0:10000 -v
```