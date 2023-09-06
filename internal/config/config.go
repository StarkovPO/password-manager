package config

import "os"

const (
	RunAddressKey           = "RUN_ADDRESS"
	DatabaseURIKey          = "DATABASE_URI"
	AccrualSystemAddressKey = "ACCRUAL_SYSTEM_ADDRESS"
	SecretKey               = "SECRET"
	SecretForPasswordKey    = "PASSWORD_SECRET"
)

// Config struct contain all data and secrets for application
type Config struct {
	RunAddressValue           string
	DatabaseURIValue          string
	AccrualSystemAddressValue string
	SecretValue               string
	PasswordSecretValue       string
}

// NewConfig create new config
func NewConfig() *Config {
	return &Config{}
}

// Init initialize config from env
func (c *Config) Init() {
	c.RunAddressValue = os.Getenv(RunAddressKey)
	c.DatabaseURIValue = os.Getenv(DatabaseURIKey)
	c.AccrualSystemAddressValue = os.Getenv(AccrualSystemAddressKey)
	c.SecretValue = os.Getenv(SecretKey)
	c.PasswordSecretValue = os.Getenv(SecretForPasswordKey)
}
