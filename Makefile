build:
	docker build -t docker-deployment -f Dockerfile .

deploy:
	docker run --rm -it --name docker-deployment \
        -e DOCKER_SERVER_IP=10.85.62.58 \
        -e DOCKER_COMPOSE_FILE=/opt/docker-compose.yml \
        -e TIMEOUT=30 \
        -e FORCE=true \
        -v ./.temp:/etc/docker/certs.d \
        -v ./config.yml:/opt/config.yml \
        -v ./docker-compose.yml:/opt/docker-compose.yml \
        docker-deployment


test:
	make build
	make deploy


