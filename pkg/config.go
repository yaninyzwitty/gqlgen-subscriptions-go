package pkg

import (
	"io"
	"log/slog"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   Server `yaml:"server"`
	Kafka    Kafka  `yaml:"kafka"`
	Database DB     `yaml:"database"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Kafka struct {
	Username         string `yaml:"username"`
	BootstrapServers string `yaml:"bootstrap_servers"`
	SecurityProtocol string `yaml:"security_protocol"`
	SASLMechanism    string `yaml:"sasl_mechanism"`
	RegistryUrl      string `yaml:"schema_registry_url"`
}

type DB struct {
	User     string `yaml:"username"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

func (c *Config) LoadConfig(file io.Reader) error {
	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("failed to read file", "error", err)
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		slog.Error("failed to unmarshal config", "error", err)
		return err
	}
	return nil

}
