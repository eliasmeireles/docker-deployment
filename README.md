Certainly! Here’s how you can incorporate the explanation about the `cert-path` and the required files into
the `README.md` file:

# Docker Container Manager

This project provides a Go script to manage Docker containers based on a configuration defined in a YAML file. The
script reads the configuration, sets up the Docker client, and runs the specified container, replacing any existing
container with the same name if the `force` flag is set to true.

## Project Structure

```
.
├── main.go
├── src
│   ├── model
│   │   └── config.go
│   ├── docker
│   │   └── client.go
│   └── container
│       └── manager.go
└── config.yml

```

## Configuration

The configuration is defined in a YAML file (`config.yml`) with the following structure:

```yaml
docker:
  host: "tcp://docker-server:2376"  # Docker daemon host URL
  tls:
    ca: "path/to/ca.pem"   # Path to the CA certificate
    cert: "path/to/cert.pem" # Path to the client certificate
    key: "path/to/key.pem"  # Path to the client key
container:
  name: "myapp_container"
  image: "myapp"
  version: "latest"
  force: true
  env:
    ENV_VAR1: "value1"
    ENV_VAR2: "value2"
  volumes:
    - "/host/path:/container/path"
  ports:
    - "8080:80"
```

## TLS Configuration

To configure Docker to use TLS, ensure your YAML file includes the following TLS settings:

| Key    | Description                    | Example Path        |
|--------|--------------------------------|---------------------|
| `ca`   | Path to the CA certificate     | `/path/to/ca.pem`   |
| `cert` | Path to the client certificate | `/path/to/cert.pem` |
| `key`  | Path to the client key         | `/path/to/key.pem`  |

### Docker Configuration

| Field  | Description                                      |
|--------|--------------------------------------------------|
| `tls`  | Whether to use TLS for Docker client connection. |
| `host` | The Docker daemon host.                          |

### Container Configuration

| Field     | Description                                                  |
|-----------|--------------------------------------------------------------|
| `name`    | The name of the container.                                   |
| `image`   | The Docker image to use.                                     |
| `version` | The version tag of the Docker image.                         |
| `force`   | Whether to replace an existing container with the same name. |
| `env`     | (Optional) Environment variables to pass to the container.   |
| `volumes` | (Optional) Volume bindings.                                  |
| `ports`   | (Optional) Port bindings.                                    |

## TLS Certificates and Keys

When using TLS for Docker, you need to provide the following files at the path specified by `cert-path`:

- **`ca.pem` (CA Certificate)**: The public certificate of the Certificate Authority (CA) that issued the server
  certificate. This is required to verify the server’s certificate.
- **`cert.pem` (Server Certificate)**: The server’s public certificate used to identify and authenticate the Docker
  server.
- **`key.pem` (Server Key)**: The server’s private key, used in conjunction with the server certificate to establish
  secure communication.
- **`ca-key.pem` (CA Key)**: The private key of the CA, used to sign the server certificate.

These files ensure encrypted communication and authentication between Docker clients and the Docker daemon.

## Usage

1. Ensure you have Docker installed and running.
2. Create your `config.yml` file based on the structure above.
3. Build and run the Go script:

    ```bash
    go build -o docker-container-manager main.go
    ./docker-container-manager config.yml
    ```

## Graceful Shutdown

The script handles graceful shutdown when receiving SIGINT or SIGTERM signals.

## Dependencies

This project uses the following dependencies:

- [docker/docker](https://github.com/docker/docker)
- [go-yaml/yaml](https://gopkg.in/yaml.v2)

To install the dependencies, run:

```bash
go get -u github.com/docker/docker/client
go get -u gopkg.in/yaml.v2
```

You can copy and paste this updated `README.md` content directly into your file. It includes details about the TLS
certificate files and their importance.
