An example of generating type-safe Envoy configuration in Go using Envoy's protocol buffers. The accompanying post explaining this repo his here: https://medium.com/@subyraman/type-safe-server-configuration-with-envoy-and-protocol-buffers-dfdf47af00d5

### Build the configuration and run the envoy server in Docker:

```
docker build --tag envoy:latest .
```

### Run the envoy server:

```
docker run -d -p 10000:10000 envoy:latest
```

Confirm that all is working by issuing a request:

```
curl 0.0.0.0:10000 -v
```
