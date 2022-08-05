FROM golang:1.18.3-alpine3.16 AS builder

RUN apk update && apk add ca-certificates && apk add tzdata
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/scraper

FROM scratch
COPY --from=builder /app/scraper /app/scraper
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8096
ENTRYPOINT ["/app/scraper"]
