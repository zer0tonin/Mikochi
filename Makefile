build:
	GOOS=linux GOARCH=amd64 go build -o mikochi github.com/zer0tonin/mikochi/src
	docker-compose build mikochi

run: build
	docker-compose up mikochi

down:
	docker-compose down --rmi=local --volumes

dev:
	docker-compose up dev -d

test:
	cd backend && go test github.com/zer0tonin/mikochi/src/...
