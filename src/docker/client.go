package docker

import (
	"docker-deployment/src/model"
	"fmt"
	"path/filepath"

	"github.com/docker/docker/client"
)

func Setup(config model.Docker) (*client.Client, error) {
	if config.Host == "" {
		return nil, fmt.Errorf("docker host is required")
	}

	if config.TLS.Ca != "" && config.TLS.Cert != "" && config.TLS.Key != "" {
		caCertPath := filepath.Join(config.TLS.Ca)
		clientCertPath := filepath.Join(config.TLS.Cert)
		clientKeyPath := filepath.Join(config.TLS.Key)

		return client.NewClientWithOpts(
			client.WithHost(config.Host),
			client.WithTLSClientConfig(caCertPath, clientCertPath, clientKeyPath),
		)
	}

	return client.NewClientWithOpts(client.WithHost(config.Host))
}
