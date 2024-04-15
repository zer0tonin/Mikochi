FROM golang:1.22

COPY ./backend /backend/
WORKDIR /backend
RUN go build -v -o mikochi .


FROM ubuntu:latest

WORKDIR /app
COPY --from=0 /backend/mikochi ./
COPY  ./frontend/dist/ ./static/

EXPOSE 8080

CMD ["./mikochi"]
