FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .
COPY src ./src/

RUN ls -lah .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deployment

# Use the official Ubuntu image
FROM eliasmeireles/dev-tools:v1

# Copy the binary from the builder stage
COPY --from=builder /app/deployment /usr/bin/deployment

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

ARG DOCKER_SERVER_IP
ARG DOCKER_COMPOSE_FILE

# Install necessary packages
RUN apt-get update && \
    apt-get install -y wget vim tar net-tools curl && \
    apt-get clean

# Download Docker binary
RUN wget https://download.docker.com/linux/static/stable/x86_64/docker-27.0.3.tgz -P /home/ubuntu/

# Extract Docker binary
RUN tar -xvf /home/ubuntu/docker-27.0.3.tgz -C /home/ubuntu/ && \
    cp /home/ubuntu/docker/* /usr/bin/ && \
    chmod +x /usr/bin/docker && \
    rm -rf /home/ubuntu/docker

# Install docker-compose
RUN curl -SL "https://github.com/docker/compose/releases/download/v2.9.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/bin/docker-compose && \
    chmod +x /usr/bin/docker-compose

# Create the directory for Docker certificates
RUN mkdir -p /etc/docker/certs.d

COPY entrypoint /usr/bin/entrypoint

# Set Docker environment variables
ENV DOCKER_COMPOSE_FILE=$DOCKER_COMPOSE_FILE
ENV DOCKER_SERVER_IP=$DOCKER_SERVER_IP
ENV DOCKER_HOST=tcp://docker-server:2376
ENV DOCKER_TLS_VERIFY=1
ENV DOCKER_CERT_PATH=/etc/docker/certs.d

# Ensure environment variables are loaded
SHELL ["/bin/bash", "-c"]

# Entry point for the container
ENTRYPOINT ["/usr/bin/entrypoint"]
