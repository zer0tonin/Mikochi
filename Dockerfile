FROM ubuntu:latest

WORKDIR /app
COPY --chmod=755 /backend/mikochi ./
COPY  ./frontend/dist/ ./static/

EXPOSE 8080

CMD ["./mikochi"]
