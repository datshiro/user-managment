export .env 

docker/up:
	docker compose  -f deploy/docker-compose.yaml up -d

docker/down:
	docker compose -f deploy/docker-compose.yaml down

run:
	[ -f .env ] && eval $(cat .env | sed 's/^/export /') || echo "no secrets file" # Export all env in .env file
	go run main.go server

db/connect:
	pgcli -h 0.0.0.0 -u postgres -W -d cake_db
