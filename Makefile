down:
	docker-compose down --rmi=local --volumes

run: down
	docker-compose up mikochi

dev: down
	docker-compose up dev -d

test:
	cd backend && go test github.com/zer0tonin/mikochi/src/...

release:
	git tag $(version)
	git push origin main --tags
	docker buildx build --platform linux/amd64,linux/arm64 -t zer0tonin/mikochi:latest -t zer0tonin/mikochi:$(version) --push .
