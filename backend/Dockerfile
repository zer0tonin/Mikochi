FROM ubuntu:latest

EXPOSE 4000

WORKDIR /app

RUN useradd mikochi
USER mikochi

COPY mikochi .

CMD ["./mikochi"]
