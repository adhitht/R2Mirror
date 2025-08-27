package config

import "fmt"

type Config struct {
	Releases []string `koanf:"releases"`
	Bucket   string   `koanf:"bucket"`
	Region   string   `koanf:"region"`
}

type StorageCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	EndpointURL     string
}

func (c *Config) Validate() error {
	if c.Bucket == "" || c.Bucket == DefaultBucketMsg {
		return fmt.Errorf("please set a valid bucket name in config.yaml")
	}
	
	if len(c.Releases) == 0 {
		return fmt.Errorf("please specify at least one Ubuntu release")
	}
	
	if c.Region == "" {
		c.Region = DefaultRegion
	}
	
	return nil
}