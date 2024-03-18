FROM golang:alpine3.19

COPY . .

RUN go mod download

RUN go build -mod=mod -o main ./cmd/service/service.go

CMD ["./main"]