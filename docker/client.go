package docker

import (
	"docker-deployment/model"
	"fmt"

	"github.com/docker/docker/client"
)

func Setup(config model.Docker) (*client.Client, error) {
	if config.Host == "" {
		return nil, fmt.Errorf("docker host is required")
	}
	if config.TLS {
		if config.CertPath == "" {
			return nil, fmt.Errorf("TLS is enabled but cert-path is not provided")
		}
		return client.NewClientWithOpts(client.WithHost(config.Host), client.WithTLSClientConfig(config.CertPath, config.CertPath, config.CertPath))
	}
	return client.NewClientWithOpts(client.WithHost(config.Host))
}
