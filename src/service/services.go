package service

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type HealthCheck struct {
	Test        []string `yaml:"test"`
	Interval    string   `yaml:"interval,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
}

type Service struct {
	ContainerName string       `yaml:"container_name"`
	Image         string       `yaml:"image"`
	Ports         []string     `yaml:"ports,omitempty"`
	DependsOn     []string     `yaml:"depends_on,omitempty"`
	HealthCheck   *HealthCheck `yaml:"healthcheck,omitempty"`
	Volumes       []string     `yaml:"volumes,omitempty"`
}

type Services struct {
	Services map[string]Service `yaml:"services"`
}

func loadServicesFromFile(filePath string) (*Services, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Decode the YAML content into the struct
	var services Services
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&services); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	return &services, nil
}
