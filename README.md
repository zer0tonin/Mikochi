# Mikochi: a minimalist remote file browser

Mikochi aims at being a remote file browser, for use in self-hosted servers / NAS.
It uses JavaScript/Preact for the frontend, and Go/Gin for the backend.
It allows you to browse remote folders, upload files, delete, rename, download and stream files to VLC/mpv.

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

Mikochi is configured using environment variabled

| Key        | Description                        |
|----------- |------------------------------------|
| host       | The ip:port mikochi will listen on |
| data_dir   | The directory accessed by mikochi  |
| jwt_secret | A secret string for jwt validation |
| username   | The username to login with         |
| password   | The password to login with         |

# Contributing

I welcome any PRs aimed at improving or fixing existing features, especially on the following subjects:

- making a linux/arm/v7 docker build (npm seems to have trouble with this architecture)
- making a smarter cache refresh
- making a fuzzier search
- improving accessibility

# Launching the development environment

The development environment and build pipeline rely on Docker/docker-compose.

Run the dockerized development environment with:
```sh
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).

Use `make build` to run a production build in a single container.
