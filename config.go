package main

import (
	"errors"
	"os"
)

type EggoConfig struct {
	AcmeEmail    string
	RedisAddress string
	APIKey       string
	ListenPort   string
	PrivateKey   string
	PublicKey    string
}

func buildConfig() (*EggoConfig, error) {
	var c *EggoConfig

	c.AcmeEmail = os.Getenv("EGGO_ACME_EMAIL")
	if c.AcmeEmail == "" {
		return nil, errors.New("EGGO_ACME_EMAIL missing")
	}
	c.RedisAddress = os.Getenv("EGGO_REDIS_ADDR")
	if c.RedisAddress == "" {
		c.RedisAddress = "localhost:6379"
	}
	c.APIKey = os.Getenv("EGGO_ZEROSSL_API_KEY")
	if c.APIKey == "" {
		return nil, errors.New("EGGO_ZEROSSL_API_KEY missing")
	}
	c.ListenPort = os.Getenv("EGGO_LISTEN_PORT")
	if c.ListenPort == "" {
		c.ListenPort = "80"
	}
	c.PrivateKey = os.Getenv("EGGO_ZEROSSL_PRIVATE_KEY")
	if c.PrivateKey == "" {
		return nil, errors.New("EGGO_ZEROSSL_PRIVATE_KEY missing")
	}
	c.PublicKey = os.Getenv("EGGO_ZEROSSL_PUBLIC_KEY")
	if c.PublicKey == "" {
		return nil, errors.New("EGGO_ZEROSSL_PUBLIC_KEY missing")
	}

	return c, nil
}
