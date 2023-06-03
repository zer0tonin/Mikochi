FROM golang:1.20

COPY ./backend /backend/
WORKDIR /backend
RUN go build -v -o mikochi .


FROM node:19.1.0-alpine

COPY ./frontend /frontend/
WORKDIR /frontend
RUN npm install
RUN npm run build


FROM ubuntu:latest

WORKDIR /app
COPY --from=0 /backend/mikochi ./
COPY --from=1 /frontend/build/ ./static/

EXPOSE 8080

CMD ["./mikochi"]
