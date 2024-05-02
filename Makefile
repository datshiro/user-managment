docker/up:
	docker compose  -f deploy/docker-compose.yaml up -d

docker/down:
	docker compose -f deploy/docker-compose.yaml down

run:
	go run main.go server -p 3000 --prefix api
