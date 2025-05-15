FROM golang:1.24.2 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .
COPY src ./src/

RUN ls -lah .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o deployment

# Use the official Ubuntu image
FROM eliasmeireles/docker-remote:latest

# Copy the binary from the builder stage
COPY --from=builder /app/deployment /usr/bin/deployment
COPY entrypoint /usr/local/bin/entrypoint
COPY start-up /usr/local/bin/start-up

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

ARG DOCKER_SERVER_IP
ARG DOCKER_COMPOSE_FILE

ENV DOCKER_COMPOSE_FILE=$DOCKER_COMPOSE_FILE

ENTRYPOINT ["/usr/local/bin/start-up"]
