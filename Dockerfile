FROM golang:alpine3.19

COPY . .

RUN go mod download

RUN go build -mod vendor -o main ./cmd/service/service.go

EXPOSE 8080

CMD ["./main"]