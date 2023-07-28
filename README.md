# 🌱 Mikochi: a minimalist remote file browser

Mikochi is a remote file browser, for use in self-hosted servers / NAS.
It allows you to browse remote folders, upload files, delete, rename, download and stream files to VLC/mpv.

It comes with a web interface powered by JavaScript/Preact, and an API built in Go/Gin.

![Screenshot of the web interface](https://github.com/zer0tonin/Mikochi/blob/main/screenshot.jpg?raw=true)

## Getting started

### Binary

Launch the app using a pre-compiled binary from the latest [release](https://github.com/zer0tonin/Mikochi/releases):

```sh
wget -c https://github.com/zer0tonin/Mikochi/releases/download/1.2.6/mikochi-linux-amd64.tar.gz -O - | tar -xz
HOST=127.0.0.1:8080 USERNAME=zer0tonin PASSWORD=horsebatterysomething ./mikochi
```

### Docker

Launch the app using [docker](https://hub.docker.com/r/zer0tonin/mikochi):

```sh
docker run \
-p 8080:8080 -v $(PWD)/data:/data \
-e DATA_DIR="/data" -e USERNAME=zer0tonin \
-e PASSWORD=horsebatterysomething zer0tonin/mikochi:latest
```

(For arm/v7 users, use the `latest-armv7` or `1.2.6-armv7` tag)

### Kubernetes

For Kubernetes users, Mikochi is installable using a [helm chart](https://artifacthub.io/packages/helm/zer0tonin/mikochi):

```sh
helm repo add zer0tonin https://zer0tonin.github.io/helm-charts/
helm install mikochi zer0tonin/mikochi \
--version 1.2.6 --set mikochi.username=zer0tonin \
--set mikochi.password=my_super_password --set persistence.enabled=true
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

## Security considerations

It is encouraged to use Mikochi behind a reverse proxy (ie. nginx), and use it to [configure TLS](nginx.org/en/docs/http/configuring_https_servers.html).
This will stop any attackers from being able to replay requests and access your files.
Additionally, using [rate limits](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html) can help slow down bruteforce and DDoS attacks.

## Contributing

I welcome any PRs aimed at improving or fixing existing features, especially on the following subjects:

- making a non-hacky linux/arm/v7 docker build (npm seems to have trouble with this architecture)
- making a smarter cache refresh
- improving accessibility
- s3/minio support

### Launching the development environment

The development environment and build pipeline rely on Docker/docker-compose.

Run the dockerized development environment with:
```sh
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).

Use `make build` to run a production build in a single container.
