build:
	docker build -t docker-deployment -f DockerfileBuilder .

deploy:
	docker run --rm -it --name docker-deployment \
        -v ./.temp:/etc/docker/certs.d \
        -v ./config.yml:/opt/config.yml \
        docker-deployment

test:
	make build
	make deploy


