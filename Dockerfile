# syntax=docker/dockerfile:1
FROM golang@sha256:d87558ea5ccfcf09ff62b08d3e5a1e1d7256499ee4d1aeb22a75da56fb5d6ac0 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN go build -o /rig ./cmd/server
RUN go install github.com/grpc-ecosystem/grpc-health-probe/v2@v0.4.39

FROM gcr.io/distroless/static-debian11@sha256:fe46af1610615bc299ea8a8e1fbe388bcf332da6eb7150110a4e81b251012c70
WORKDIR /
COPY --from=builder /rig /rig
COPY --from=builder /go/bin/grpc-health-probe /usr/local/bin/grpc-health-probe
EXPOSE 50051
HEALTHCHECK CMD ["/usr/local/bin/grpc-health-probe", "-addr=:50051"]
ENTRYPOINT ["/rig"]
