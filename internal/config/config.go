package config

import "os"

const (
	RunAddressKey           = "RUN_ADDRESS"
	DatabaseURIKey          = "DATABASE_URI"
	AccrualSystemAddressKey = "ACCRUAL_SYSTEM_ADDRESS"
)

type Config struct {
	RunAddressValue           string
	DatabaseURIValue          string
	AccrualSystemAddressValue string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Init() {
	c.RunAddressValue = os.Getenv(RunAddressKey)
	c.DatabaseURIValue = os.Getenv(DatabaseURIKey)
	c.AccrualSystemAddressValue = os.Getenv(AccrualSystemAddressKey)
}
