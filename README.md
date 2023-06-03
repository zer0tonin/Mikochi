# 🌱 Mikochi: a minimalist remote file browser

Mikochi is a remote file browser, for use in self-hosted servers / NAS.
It allows you to browse remote folders, upload files, delete, rename, download and stream files to VLC/mpv.

It comes with a web interface powered by JavaScript/Preact, and an API built in Go/Gin.

![Screenshot of the web interface](https://github.com/zer0tonin/Mikochi/blob/main/screenshot.jpg?raw=true)

# Getting started

## Binary

## Docker

Launch the app using docker:

```sh
docker run -v ~/Code/Mikochi/data:/data \
-p 8080:8080 -e host="0.0.0.0:8080" \
-e data_dir="/data" -e jwt_secret=my_secret \
-e username=root -e password=pass \
zer0tonin/mikochi:latest
```

## Kubernetes


# Configuration

Mikochi is configured using environment variables

| Key        | Description                        | Default    |
|----------- |------------------------------------|------------|
| HOST       | The ip:port mikochi will listen on | 0.0.0.0:80 |
| DATA_DIR   | The directory accessed by mikochi  | /data      |
| JWT_SECRET | A secret string for jwt validation | [Random]   |
| USERNAME   | The username to login with         | root       |
| PASSWORD   | The password to login with         | pass       |

# Security considerations

It is encouraged to use Mikochi behind a reverse proxy (ie. nginx), and use it to [configure TLS](nginx.org/en/docs/http/configuring_https_servers.html).
This will stop any attackers from being able to replay requests and access your files.
Additionally, using [rate limits](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html) can help slow down bruteforce and DDoS attacks.

# Contributing

I welcome any PRs aimed at improving or fixing existing features, especially on the following subjects:

- making a linux/arm/v7 docker build (npm seems to have trouble with this architecture)
- making a smarter cache refresh
- making a fuzzier search
- improving accessibility
- s3/minio support

## Launching the development environment

The development environment and build pipeline rely on Docker/docker-compose.

Run the dockerized development environment with:
```sh
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).

Use `make build` to run a production build in a single container.
