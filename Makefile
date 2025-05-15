# Makefile for building and pushing a Docker image

.PHONY: build push buildx all

# Setup Buildx builder
buildx:
	@docker buildx create --name buildxBuilder --use
	@docker buildx inspect buildxBuilder --bootstrap

# Build the Docker image
build:
	@echo "Building Docker image..."
	@read -p "Enter the tag version: " TAG; \
	 docker buildx build --no-cache --platform linux/amd64,linux/arm64 -t eliasmeireles/docker-deployment:$$TAG --push .

image-build:
	@echo "Building Docker image..."
	@read -p "Enter the tag version: " TAG; \
	 docker build -t eliasmeireles/docker-deployment:$$TAG .

# Push the Docker image
push:
	@echo "Pushing Docker image..."
	@read -p "Enter the tag version: " TAG; \
	 docker push eliasmeireles/docker-deployment:$$TAG

# Default target
all: build

local-build:
	go build -o docker-deployment main.go
	chmod +x  docker-deployment

test:
	make local-build
	DOCKER_COMPOSE_FILE=./example/docker-compose.yml FORCE=true TIMEOUT=300 ./docker-deployment

container-run:
	 docker run --rm -it \
	   --name docker-deployment \
	   -v ./example:/opt \
	   -v ~/docker/certs.d:/etc/docker/certs.d \
	   -e DOCKER_COMPOSE_FILE=/opt/docker-compose.yml \
	   -e DOCKER_REGISTRY_HOST=$$DOCKER_REGISTRY_HOST \
	   -e DOCKER_REGISTRY_PASSWORD=$$DOCKER_REGISTRY_PASSWORD \
	   -e DOCKER_SERVER_IP=$$DOCKER_SERVER_IP \
	   -e DOCKER_REGISTRY_USERNAME=$$DOCKER_REGISTRY_USERNAME \
	   eliasmeireles/docker-deployment:v1



