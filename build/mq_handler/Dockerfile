FROM golang:1.19.2-alpine3.16 AS builder
WORKDIR /app
COPY go.mod .
COPY /go.sum .
COPY config/mq_handler/config-docker.yaml .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/mqhandler /cmd/mq_handler/main.go

FROM scratch
COPY --from=builder /app/config/mq_handler/config-docker.yaml /app/config/mq_handler/config-docker.yaml
COPY --from=builder /app/mqhandler /app/mqhandler
EXPOSE 8091
ENTRYPOINT ["/app/mqhandler"]
