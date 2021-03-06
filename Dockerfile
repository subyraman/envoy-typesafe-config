FROM golang:1.13.0-stretch AS builder

WORKDIR /go/src/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
RUN go build main.go

RUN ./main


FROM envoyproxy/envoy-dev
RUN apt-get update && apt-get install -y netcat dnsutils
#COPY  --from=builder /go/src/envoy.json /etc/envoy/envoy.json

COPY envoy/envoy.yaml /etc/envoy/envoy.yaml

EXPOSE 10000
#CMD /usr/local/bin/envoy -c /etc/envoy/envoy.json
CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml