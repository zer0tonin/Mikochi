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

release:
	git tag $(version)
	git push origin main --tags
	docker buildx build --platform linux/amd64,linux/arm64 -t zer0tonin/mikochi-frontend:latest -t zer0tonin/mikochi-frontend:$(version) --push -f ./frontend/Dockerfile-prod ./frontend
	docker buildx build --platform linux/amd64,linux/arm64 -t zer0tonin/mikochi-frontend:latest -t zer0tonin/mikochi-backend:$(version) --push ./backend
