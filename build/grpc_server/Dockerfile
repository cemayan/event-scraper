FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/grpc_server/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/grpc-server /cmd/grpc_server/main.go

FROM scratch
COPY --from=builder /app/config/grpc_server/config-docker.yaml /app/config/grpc_server/config-docker.yaml
COPY --from=builder /app/grpc-server /app/grpc-server
EXPOSE 50051
ENTRYPOINT ["/app/grpc-server"]
