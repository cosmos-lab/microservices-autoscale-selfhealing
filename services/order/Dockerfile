FROM docker.io/library/golang:1.24-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go clean -modcache
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o app

FROM docker.io/library/alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
