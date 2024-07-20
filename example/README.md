# Docker Remote Access Setup

This setup automates the installation and configuration of Docker with TLS on a Multipass instance (`docker-server`) and
prepares another Multipass instance (`docker-connector`) to access the Docker server remotely.

## Prerequisites

- **Multipass**: Ensure that Multipass is installed on your system.
- **Bash**: A Unix-like environment to run shell scripts.

## Setup Overview

1. **Docker Remote Access Configuration Script**:
    - Installs Docker on the `docker-server`.
    - Generates TLS certificates for secure Docker communication.
    - Configures Docker to use the generated TLS certificates.

2. **Multipass Setup Script**:
    - Deletes any existing `docker-server` and `docker-connector` instances.
    - Creates new `docker-server` and `docker-connector` instances.
    - Copies and runs the Docker remote access configuration script on `docker-server`.
    - Sets up `docker-connector` to access the `docker-server` securely via TLS.

## Steps to Run the Setup

1. **Clone the Repository**:
    ```sh
    git clone git@github.com:eliasmeireles/docker-deployment.git
    ```

2. **Navigate to the `docker-deployment/example` Directory**:
    ```sh
    cd docker-deployment/example
    ./config-test  
    ```


