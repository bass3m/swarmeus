package config

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type Target struct {
	Job           string   `yaml:"job"`
	InstanceRegex string   `yaml:"instance_regex"`
	Port          int64    `yaml:"port"`
	MetricsPath   string   `yaml:"metrics_path"`
	Labels        []string `yaml:"labels,flow,omitempty"`
}

type Config struct {
	Swarmeus struct {
		ScanInterval time.Duration `yaml:"scan_interval"`
		Network      string        `yaml:"network"`
		Endpoint     string        `yaml:"endpoint"`
		DockerMode   string        `yaml:"docker_mode"`
		SDFilePath   string        `yaml:"sd_file_path"`
	}
	Targets []Target `yaml:"targets,flow"`
}

func ReadConfig(configPath string) (Config, error) {
	var cfg Config

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Errorf("Failed to read YAML config file err:  %v ", err)
		return Config{}, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
		return Config{}, err
	}

	return cfg, nil
}
