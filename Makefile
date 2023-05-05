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


release: build
	git tag $(version)
	git push origin main --tags
	docker tag mikochi-frontend:latest zer0tonin/mikochi-frontend:latest
	docker push zer0tonin/mikochi-frontend:latest
	docker tag mikochi-frontend:latest zer0tonin/mikochi-frontend:$(version)
	docker push zer0tonin/mikochi-frontend:$(version)
	docker tag mikochi-backend:latest zer0tonin/mikochi-backend:latest
	docker push zer0tonin/mikochi-backend:latest
	docker tag mikochi-backend:latest zer0tonin/mikochi-backend:$(version)
	docker push zer0tonin/mikochi-backend:$(version)
