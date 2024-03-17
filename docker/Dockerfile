FROM golang:alpine

COPY . .

RUN go build -o main ./cmd/service/service.go

CMD ["./main"]