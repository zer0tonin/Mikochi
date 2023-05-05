FROM golang:1.20

EXPOSE 4000

WORKDIR /app

COPY . .

RUN go build -v -o mikochi ./...

CMD ["./mikochi"]
