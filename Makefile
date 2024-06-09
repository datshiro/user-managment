export .env 

up:
	docker compose --env-file=.env -f deploy/docker-compose.yaml up -d

down:
	docker compose -f deploy/docker-compose.yaml down

build:
	docker compose -f deploy/docker-compose.yaml build

logs:
	docker logs -f my-app

db/connect:
	pgcli -h 0.0.0.0 -u postgres -W -d my_db

test:
	go test -v ./... 



dind:
	docker network create jenkins || true
	docker run \
  --name jenkins-docker \
  --rm \
  --detach \
  --privileged \
  --network jenkins \
  --network-alias docker \
  --env DOCKER_TLS_CERTDIR=/certs \
  --volume jenkins-docker-certs:/certs/client \
  --volume jenkins-data:/var/jenkins_home \
  --publish 2376:2376 \
  docker:dind \
  --storage-driver overlay2

jenkins:
	docker run \
		--name jenkins-instance \
		--restart=on-failure \
		--detach \
		--network jenkins \
		--env DOCKER_HOST=tcp://docker:2376 \
		--env DOCKER_CERT_PATH=/certs/client \
		--env DOCKER_TLS_VERIFY=1 \
		--publish 8080:8080 \
		--publish 50000:50000 \
		--volume jenkins-data:/var/jenkins_home \
		--volume jenkins-docker-certs:/certs/client:ro \
		jenkins/jenkins
