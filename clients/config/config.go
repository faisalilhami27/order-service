package config

import "github.com/parnurzeal/gorequest"

type IClientConfig interface {
	Client() *gorequest.SuperAgent
	BaseURL() string
	SecretKey() string
}

type ClientConfig struct {
	client    *gorequest.SuperAgent
	baseURL   string
	secretKey string
}

type Option func(*ClientConfig)

func NewClientConfig(options ...Option) IClientConfig {
	clientConfig := &ClientConfig{
		client: gorequest.New().
			Set("Content-Type", "application/json").
			Set("Accept", "application/json"),
	}
	for _, option := range options {
		option(clientConfig)
	}
	return clientConfig
}

func (c *ClientConfig) Client() *gorequest.SuperAgent {
	return c.client
}

func (c *ClientConfig) BaseURL() string {
	return c.baseURL
}

func (c *ClientConfig) SecretKey() string {
	return c.secretKey
}

func WithBaseURL(baseURL string) Option {
	return func(c *ClientConfig) {
		c.baseURL = baseURL
	}
}

func WithSecretKey(secretKey string) Option {
	return func(c *ClientConfig) {
		c.secretKey = secretKey
	}
}
