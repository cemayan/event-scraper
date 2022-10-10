FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/auth/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o/app/auth-server cmd/auth/main.go


FROM scratch
COPY --from=builder /app/config/auth/config-docker.yaml /app/config/auth/config-docker.yaml
COPY --from=builder /app/auth-server /app/auth-server
EXPOSE 8109
ENTRYPOINT ["/app/auth-server"]
