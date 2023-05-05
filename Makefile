down:
	docker-compose down --rmi=local --volumes

build: down
	docker-compose build mikochi

run: build
	docker-compose up mikochi -d

dev: down
	docker-compose up dev -d

test:
	cd backend && go test github.com/zer0tonin/mikochi/src/...
