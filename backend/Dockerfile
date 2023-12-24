FROM golang:1.21

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN apt-get update
RUN apt-get install npm -y
RUN npm install -g nodemon

WORKDIR /mikochi

CMD ["nodemon", "--exec", "dlv debug --continue --accept-multiclient --listen=:2345 --headless=true --api-version=2 --log .", "-e", "go,mod,sum", "--signal", "SIGINT", "--watch", "/mikochi"]
