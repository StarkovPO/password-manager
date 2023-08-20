package config

import "os"

const (
	RunAddressKey           = "RUN_ADDRESS"
	DatabaseURIKey          = "DATABASE_URI"
	AccrualSystemAddressKey = "ACCRUAL_SYSTEM_ADDRESS"
	SecretKey               = "SECRET"
	SecretForPasswordKey    = "PASSWORD_SECRET"
)

type Config struct {
	RunAddressValue           string
	DatabaseURIValue          string
	AccrualSystemAddressValue string
	SecretValue               string
	PasswordSecretValue       string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() {
	c.RunAddressValue = os.Getenv(RunAddressKey)
	c.DatabaseURIValue = os.Getenv(DatabaseURIKey)
	c.AccrualSystemAddressValue = os.Getenv(AccrualSystemAddressKey)
	c.SecretValue = os.Getenv(SecretKey)
	c.PasswordSecretValue = os.Getenv(SecretForPasswordKey)
}
