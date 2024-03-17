FROM golang:alpine

COPY . .

RUN go build -o main ./service/service.go

CMD ["./main"]