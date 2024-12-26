package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the service configuration
type Config struct {
	App struct {
		Name        string `yaml:"name"`
		Environment string `yaml:"environment"`
	} `yaml:"app"`

	Server struct {
		Port    int `yaml:"port"`
		Timeout struct {
			Read  time.Duration `yaml:"read"`
			Write time.Duration `yaml:"write"`
			Idle  time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`

	Database struct {
		Host            string        `yaml:"host"`
		Port            int           `yaml:"port"`
		User            string        `yaml:"user"`
		Password        string        `yaml:"password"`
		Name            string        `yaml:"name"`
		MaxOpenConns    int           `yaml:"max_open_conns"`
		MaxIdleConns    int           `yaml:"max_idle_conns"`
		ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	} `yaml:"database"`

	NATS struct {
		URL           string `yaml:"url"`
		ClusterID     string `yaml:"cluster_id"`
		ClientID      string `yaml:"client_id"`
		SubjectPrefix string `yaml:"subject_prefix"`
	} `yaml:"nats"`

	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
		Output string `yaml:"output"`
	} `yaml:"logging"`

	Metrics struct {
		Enabled        bool `yaml:"enabled"`
		PrometheusPort int  `yaml:"prometheus_port"`
	} `yaml:"metrics"`
}

// Load reads the configuration from a YAML file
func Load() (*Config, error) {
	config := &Config{}

	// Look for config file in different locations
	configPaths := []string{
		"configs/config.yaml",
		"../configs/config.yaml",
		"../../configs/config.yaml",
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		return nil, fmt.Errorf("config file not found in any of the search paths")
	}

	// Read the config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	// Override with environment variables if present
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		var port int
		if _, err := fmt.Sscanf(envPort, "%d", &port); err == nil {
			config.Server.Port = port
		}
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		var port int
		if _, err := fmt.Sscanf(dbPort, "%d", &port); err == nil {
			config.Database.Port = port
		}
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}

	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}

	if natsURL := os.Getenv("NATS_URL"); natsURL != "" {
		config.NATS.URL = natsURL
	}

	// Validate required fields
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation error: %v", err)
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.App.Name == "" {
		return fmt.Errorf("app name is required")
	}

	if config.Server.Port == 0 {
		return fmt.Errorf("server port is required")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port == 0 {
		return fmt.Errorf("database port is required")
	}

	if config.NATS.URL == "" {
		return fmt.Errorf("NATS URL is required")
	}

	return nil
}
