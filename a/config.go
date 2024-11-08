package a

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Profile struct {
	SrcDir  string `yaml:"srcDir"`
	DestDir string `yaml:"destDir"`
}

type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

const configFilePath = "./config/config.yaml"

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// 返回默认值
			return &Config{
				Profiles: map[string]Profile{
					"photos": {SrcDir: "/sdcard/Download", DestDir: "/Users/jhonhe/Photos"},
					"files":  {SrcDir: "/sdcard/Download", DestDir: "/Users/jhonhe/Files"},
				},
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
