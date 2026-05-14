# 🌱 Mikochi: a minimalist remote file browser

Mikochi is a remote file browser for your self-hosted server or NAS.
It allows you to browse remote folders, upload, delete, rename and download files.
You can also use it to generate streaming links which can be played using VLC or mpv

<img width="888" height="480" alt="mikochi" src="https://github.com/user-attachments/assets/48141ba4-7015-45fa-81ee-dfecb347c48f" />


## Features

* Browse files on you remote server
* Fuzzy search
* Upload new files and create folders
* Rename and delete files
* Download files and directories (in .tar.gz)
* Stream files to VLC/MPV
* Lightweight web interface built in Preact
* High-performance server built in Go/Gin

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
--version 1.10.0 --set mikochi.username=zer0tonin \
--set mikochi.password=my_super_password --set persistence.enabled=true
```

### Debian / Ubuntu / Mint

Install the app using a .deb package from the latest [release](https://github.com/zer0tonin/Mikochi/releases).

```sh
wget -c https://github.com/zer0tonin/Mikochi/releases/download/1.10.0/mikochi-1.10.0-linux-amd64.deb

sudo chmod +x mikochi-1.10.0-linux-amd64.deb
sudo apt install mikochi-1.10.0-linux-amd64.deb
sudo mkdir /data

PASSWORD=horsebatterysomething mikochi
```

### Fedora / Red Hat / CentOS

Install the app using a .rpm package from the latest [release](https://github.com/zer0tonin/Mikochi/releases).

```sh
wget -c https://github.com/zer0tonin/Mikochi/releases/download/1.10.0/mikochi-1.10.0-linux-amd64.rpm

sudo chmod +x mikochi-1.10.0-linux-amd64.rpm
sudo rpm -ivh mikochi-1.10.0-linux-amd64.rpm
sudo mkdir /data

PASSWORD=horsebatterysomething mikochi
```

### Binary

Launch the app using a pre-compiled binary from the latest [release](https://github.com/zer0tonin/Mikochi/releases):

```sh
wget -c https://github.com/zer0tonin/Mikochi/releases/download/1.10.0/mikochi-1.10.0-linux-amd64.tar.gz -O - | tar -xz

sudo mkdir /data
sudo mkdir /usr/share/mikochi
sudo mv ./mikochi-1.10.0/mikochi /usr/bin/mikochi
sudo mv ./mikochi-1.10.0/static /usr/share/mikochi/static

PASSWORD=horsebatterysomething mikochi
```

## Configuration

Mikochi is configured using environment variables

| Key            | Description                                 | Default                     |
|----------------|---------------------------------------------|-----------------------------|
| HOST           | The ip:port mikochi will listen on          | 0.0.0.0:80                  |
| DATA_DIR       | The directory accessed by mikochi           | /data                       |
| JWT_SECRET     | A secret string for jwt validation          | [Random]                    |
| USERNAME       | The username to login with                  | root                        |
| PASSWORD       | The password to login with                  | pass                        |
| CERT_CA        | The path to a TLS certificate               | null                        |
| CERT_KEY       | The path to the key associated with CERT_CA | null                        |
| NO_AUTH        | If true, disables all authentication        | false                       |
| GZIP           | If true, enables gzip compression           | false                       |
| FRONTEND_FILES | The location of the frontend static files   | /usr/share/mikochi/static   |

Note: it is recommended to not manually set JWT_SECRET, as getting a new randomly generated secret everytime when mikochi starts let's you invalidate authentication tokens by restarting the process.


## Launching the development environment

The development environment and build pipeline rely on Docker/docker-compose.

Run the dockerized development environment with:
```sh
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).
