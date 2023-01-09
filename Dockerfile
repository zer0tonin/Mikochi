FROM ubuntu:latest

EXPOSE 8080

WORKDIR /app

RUN useradd mikochi
USER mikochi

COPY mikochi .
COPY templates ./templates
COPY config.yaml .

CMD ["./mikochi"]
