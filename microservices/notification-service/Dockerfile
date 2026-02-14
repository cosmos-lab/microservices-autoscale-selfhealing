FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-service ./main.go

EXPOSE 8083

CMD ["./notification-service"]

