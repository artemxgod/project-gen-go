package configs

	import (
			"os"
	
			"gopkg.in/yaml.v3"
	)
	
	type Config struct {
			ServiceName string 
	}

	func ReadConfig(path string) (*Config, error) {
			configFile, err := os.ReadFile(path)
			if err != nil {
					return nil, err
			}
	
			return readConfigFromFile(configFile)
	}
	
	func readConfigFromFile(fileBytes []byte) (*Config, error) {
			cfg := new(Config)
			if err := yaml.Unmarshal(fileBytes, cfg); err != nil {
					return nil, err
			}
	
			return cfg, nil
	}