down:
	ls backend | grep __debug_bin | gxargs -I {} rm backend/{}
	docker-compose down --rmi=local --volumes

run: down
	docker-compose up mikochi -d

dev: down
	docker-compose up dev -d

test:
	cd backend && go test github.com/zer0tonin/mikochi/src/...
