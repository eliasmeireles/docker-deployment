package model

type Docker struct {
	TLS      bool   `yaml:"tls"`
	CertPath string `yaml:"cert-path"`
	Host     string `yaml:"host"`
}

type Container struct {
	Name    string            `yaml:"name"`
	Image   string            `yaml:"image"`
	Version string            `yaml:"version"`
	Force   bool              `yaml:"force"`
	Env     map[string]string `yaml:"env"`
	Volumes []string          `yaml:"volumes"`
	Port    []string          `yaml:"ports"`
}

type Config struct {
	DockerConfig    Docker    `yaml:"docker"`
	ContainerConfig Container `yaml:"container"`
}
