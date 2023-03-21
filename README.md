# Mikochi: a minimalist remote file browser

Mikochi aims at being a remote file browser, for use in self-hosted servers / NAS.
It uses JavaScript/Preact for the frontend, and Go/Gin for the backend.
This is currently a *work in progress*.

It has 3 main features:
* browse files
* search files
* stream/download a file

Planned features:
* delete a file
* rename a file

# How to run it

Run the dockerized development environment with:
```
make dev
```

It will start a frontend container (listening on 5000), a backend container (listening on 4000) and an nginx to wire both (listening on 8080).

Create a production build with:
```
make build
```
