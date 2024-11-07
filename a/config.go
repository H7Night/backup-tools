package a

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SrcDir  string `yaml:"srcDir"`
	DestDir string `yaml:"destDir"`
}

const configFilePath = "./config/config.yaml"

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				SrcDir:  "/sdcard/Download",
				DestDir: "/User/jhonhe/Downloads",
			}, nil
		}
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, data, 0644)
}
