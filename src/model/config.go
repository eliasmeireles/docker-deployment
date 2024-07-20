package model

type TLS struct {
	Ca   string `yaml:"ca"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Docker struct {
	TLS  TLS    `yaml:"tls"`
	Host string `yaml:"host"`
}

type Startup struct {
	Timeout int `yaml:"timeout"`
}

type HealthCheck struct {
	Test        []string `yaml:"test"`         // Command or command arguments for health check
	Interval    int      `yaml:"interval"`     // Interval between health checks in seconds
	Timeout     int      `yaml:"timeout"`      // Timeout for each health check in seconds
	Retries     int      `yaml:"retries"`      // Number of retries before marking unhealthy
	StartPeriod int      `yaml:"start_period"` // Delay before starting health checks in seconds
}

type Container struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Version     string            `yaml:"version"`
	Force       bool              `yaml:"force"`
	Env         map[string]string `yaml:"env"`
	Volumes     []string          `yaml:"volumes"`
	Ports       []string          `yaml:"ports"`
	Startup     Startup           `yaml:"startup"`
	HealthCheck *HealthCheck      `yaml:"healthcheck"` // Optional: Health check configuration
}

type Config struct {
	DockerConfig    Docker    `yaml:"docker"`
	ContainerConfig Container `yaml:"container"`
}
