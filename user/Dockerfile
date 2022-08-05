FROM golang:1.18.3-alpine3.16 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/user-service

FROM scratch
COPY --from=builder /app/user-service /app/user-service
EXPOSE 8089
ENTRYPOINT ["/app/user-service"]
