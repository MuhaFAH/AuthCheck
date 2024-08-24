FROM golang:1.22-alpine

WORKDIR /

COPY . .

RUN go mod download

RUN go build -o main ./cmd/main.go

CMD ["./main"]