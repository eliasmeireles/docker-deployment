# Define image version and container name
IMAGE_NAME=eliasmeireles/docker-deployment
IMAGE_VERSION=v1
CONTAINER_NAME=docker-deployment

# Define environment variables
DOCKER_SERVER_IP=docker-server # IP address of the Docker server
DOCKER_COMPOSE_FILE=/opt/docker-compose.yml
TIMEOUT=30
FORCE=true

# Define paths for volumes
COMPOSE_FILE_PATH=./docker-compose.yml
CERTS_VOLUME=/etc/docker/certs.d
HOST_CERTS=/etc/docker/certs.d
COMPOSE_VOLUME=/opt/docker-compose.yml
COMPOSE_VOLUME=/opt/docker-compose.yml

# Define deploy target
deploy:
	docker run --rm -it --name $(CONTAINER_NAME) \
		-e DOCKER_SERVER_IP=$(DOCKER_SERVER_IP) \
		-e DOCKER_COMPOSE_FILE=$(DOCKER_COMPOSE_FILE) \
		-e TIMEOUT=$(TIMEOUT) \
		-e FORCE=$(FORCE) \
		-v $(HOST_CERTS):$(CERTS_VOLUME) \
		-v ./example/test:/opt \
		-v $(COMPOSE_FILE_PATH):$(COMPOSE_VOLUME) \
		$(IMAGE_NAME):$(IMAGE_VERSION)

copy-connector:
	multipass transfer docker-compose.yml docker-connector:/home/ubuntu/docker-compose.yml
	multipass transfer Makefile docker-connector:/home/ubuntu/Makefile

server:
	./config-test
