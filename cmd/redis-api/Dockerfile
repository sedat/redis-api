FROM golang:1.17

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/redis-api
EXPOSE 9000
CMD ["./main"]

