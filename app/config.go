package app

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"time"
)

const defaultCallTimout = 10 * time.Second

type Config struct {
	kubeClientConfig clientcmd.ClientConfig
	callTimout       time.Duration
}

// NewConfig returns a new k8s client config
func NewConfig() *Config {
	// Create a new rest client
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	return &Config{
		kubeClientConfig: config,
		callTimout:       defaultCallTimout,
	}
}

func (c *Config) RESTClient() (*rest.Config, error) {
	return c.kubeClientConfig.ClientConfig()
}

func (c *Config) RawConfig() (api.Config, error) {
	return c.kubeClientConfig.RawConfig()
}

func (c *Config) ConfigAccess() clientcmd.ConfigAccess {
	return c.kubeClientConfig.ConfigAccess()
}
