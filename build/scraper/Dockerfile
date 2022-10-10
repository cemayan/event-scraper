FROM golang:1.19.2-alpine3.16 AS builder
RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY config/scraper/config-docker.yaml .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/scraper  cmd/scraper/main.go

FROM scratch
COPY --from=builder /app/config/scraper/config-docker.yaml /app/config/scraper/config-docker.yaml
COPY --from=builder /app/scraper /app/scraper
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8096
ENTRYPOINT ["/app/scraper"]
