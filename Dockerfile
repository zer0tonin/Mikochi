FROM ubuntu:latest

WORKDIR /app
COPY /backend/mikochi ./
COPY  ./frontend/dist/ ./static/

EXPOSE 8080

CMD ["./mikochi"]
