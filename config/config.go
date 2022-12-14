package config

import "github.com/caarlos0/env/v6"

// Config contains the env variables needed to run the servers
type Config struct {
	DBHost            string `env:"POSTGRES_HOST,required"`
	DBPort            int    `env:"POSTGRES_PORT,required"`
	DBUser            string `env:"POSTGRES_USER,required"`
	DBPassword        string `env:"POSTGRES_PASSWORD,required"`
	DBName            string `env:"POSTGRES_NAME,required"`
	SigningKey        string `env:"SIGNING_KEY,required"`
	GoogleCredentials string `env:"GOOGLE_CREDENTIALS,required"`
	SendGridAPIKey    string `env:"SENDGRID_API_KEY"`
	Domain            string `env:"DOMAIN,required"`
	RedirectURL       string `env:"REDIRECT_URL,required"`
}

func NewFromEnv() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}
