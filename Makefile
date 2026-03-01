down:
	ls backend | grep __debug_bin | xargs -I {} rm backend/{}
	docker compose down --rmi=local --volumes

run: down
	cd backend && GOOS=linux go build -v -o ./mikochi .
	cd frontend && npm run build
	docker compose build mikochi
	docker compose up mikochi

dev: down
	docker compose up dev -d

test:
	cd backend && go test github.com/zer0tonin/mikochi/src/...
