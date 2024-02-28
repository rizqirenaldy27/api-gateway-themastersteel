FROM golang:1.21.4

WORKDIR /app

COPY . .

RUN go build -o service-user

CMD ["./service-user"]
