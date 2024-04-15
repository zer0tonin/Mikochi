# 🌱 Mikochi: a minimalist remote file browser

Mikochi is a remote file browser, for use in self-hosted servers / NAS.
It allows you to browse remote folders, upload files, delete, rename, download and stream files to VLC/mpv.

It comes with a web interface powered by JavaScript/Preact, and an API built in Go/Gin.

![Screenshot of the web interface](https://github.com/zer0tonin/Mikochi/blob/main/screenshot.jpg?raw=true)

## Getting started

### Docker

Launch the app using [docker](https://hub.docker.com/r/zer0tonin/mikochi):

```sh
docker run \
-p 8080:8080 -v $(PWD)/data:/data \
-e DATA_DIR="/data" -e USERNAME=zer0tonin \
-e PASSWORD=horsebatterysomething zer0tonin/mikochi:latest
```

You'll find a complete tutorial on installing Mikochi securely with Docker and Traefik [here](https://alicegg.tech/2024/01/04/mikochi-tutorial).

### Kubernetes

For Kubernetes users, Mikochi is installable using a [helm chart](https://artifacthub.io/packages/helm/zer0tonin/mikochi):

```sh
helm repo add zer0tonin https://zer0tonin.github.io/helm-charts/
helm install mikochi zer0tonin/mikochi \
--version 1.4.5 --set mikochi.username=zer0tonin \
--set mikochi.password=my_super_password --set persistence.enabled=true
```

### Binary

Launch the app using a pre-compiled binary from the latest [release](https://github.com/zer0tonin/Mikochi/releases):

```sh
wget -c https://github.com/zer0tonin/Mikochi/releases/download/1.4.5/mikochi-linux-amd64.tar.gz -O - | tar -xz
HOST=127.0.0.1:8080 USERNAME=zer0tonin PASSWORD=horsebatterysomething ./mikochi
```

## Configuration

Mikochi is configured using environment variables

| Key        | Description                        | Default    |
|----------- |------------------------------------|------------|
| HOST       | The ip:port mikochi will listen on | 0.0.0.0:80 |
| DATA_DIR   | The directory accessed by mikochi  | /data      |
| JWT_SECRET | A secret string for jwt validation | [Random]   |
| USERNAME   | The username to login with         | root       |
| PASSWORD   | The password to login with         | pass       |


## Launching the development environment

The development environment and build pipeline rely on Docker/docker-compose.

Run the dockerized development environment with:
```sh
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).

Use `make build` to run a production build in a single container.
