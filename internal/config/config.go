package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/adhitht/R2Mirror/pkg"
)

const (
	ConfigFile = "config.yaml"
	EnvFile    = ".env"
	
	EnvR2AccessKeyID     = "R2_ACCESS_KEY_ID"
	EnvR2SecretAccessKey = "R2_SECRET_ACCESS_KEY"
	EnvR2EndpointURL     = "R2_ENDPOINT_URL"
	EnvR2DefaultRegion   = "R2_DEFAULT_REGION"
	
	EnvAWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	EnvAWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	EnvAWSEndpointURL     = "AWS_ENDPOINT_URL"
	
	DefaultRegion    = "auto"
	DebounceDelay    = 300 * time.Millisecond
	DefaultBucketMsg = "your-bucket-name"
)

const DefaultConfigYAML = `# Ubuntu Release Downloader Configuration
releases:
  - "22.04"
  - "20.04"
  - "18.04"

# R2/S3 Configuration
bucket: "your-bucket-name"
region: "auto"  # Use "auto" for Cloudflare R2
`

const DefaultEnvContent = `# Cloudflare R2 Configuration
R2_ACCESS_KEY_ID=your-access-key-here
R2_SECRET_ACCESS_KEY=your-secret-key-here
R2_ENDPOINT_URL=https://your-account-id.r2.cloudflarestorage.com

# Optional: Set default region
# R2_DEFAULT_REGION=auto

# Fallback AWS compatibility
# AWS_ACCESS_KEY_ID=your-access-key-here
# AWS_SECRET_ACCESS_KEY=your-secret-key-here
`

var k = koanf.New(".")

func LoadEnv() error {
	if !pkg.FileExists(EnvFile) {
		fmt.Println("Creating default .env")
		if err := pkg.WriteFile(EnvFile, DefaultEnvContent); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}
		fmt.Println("Default .env created. Please edit with your R2 credentials.")
		return nil
	}

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	fmt.Println("Loaded .env")
	return nil
}

func Load() (*Config, error) {
	if err := createDefaultConfig(); err != nil {
		return nil, err
	}

	if err := k.Load(file.Provider(ConfigFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("failed to load config.yaml: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	fmt.Printf("📋 Config loaded: %d releases, bucket: %s\n", len(cfg.Releases), cfg.Bucket)
	return &cfg, nil
}

func GetStorageCredentials() (*StorageCredentials, error) {
	accessKey := pkg.GetEnvWithFallback(EnvR2AccessKeyID, EnvAWSAccessKeyID)
	secretKey := pkg.GetEnvWithFallback(EnvR2SecretAccessKey, EnvAWSSecretAccessKey)
	endpointURL := pkg.GetEnvWithFallback(EnvR2EndpointURL, EnvAWSEndpointURL)

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("storage credentials not found. Set %s and %s in .env", 
			EnvR2AccessKeyID, EnvR2SecretAccessKey)
	}

	return &StorageCredentials{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		EndpointURL:     endpointURL,
	}, nil
}

func createDefaultConfig() error {
	if !pkg.FileExists(ConfigFile) {
		fmt.Println("📝 Creating default config.yaml...")
		if err := pkg.WriteFile(ConfigFile, DefaultConfigYAML); err != nil {
			return fmt.Errorf("failed to create config.yaml: %w", err)
		}
		fmt.Println("Default config.yaml created. Please edit with your settings.")
	}
	return nil
}