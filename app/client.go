package app

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientProvider holds kube client and its config
type ClientProvider struct {
	kubeClient kubernetes.Interface
	config     *Config
}

// NewClient creates a k8s client
func NewClient(config *Config) (*ClientProvider, error) {
	restClient, err := config.RESTClient()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(restClient)
	if err != nil {
		return nil, err
	}
	return &ClientProvider{
		kubeClient: client,
		config:     config,
	}, nil
}

// CheckServerConnection checks connection to k8s server
func (c *ClientProvider) CheckServerConnection() (bool, error) {
	if c.kubeClient == nil {
		return false, fmt.Errorf("client is not initiliazed")
	}
	_, err := c.kubeClient.Discovery().ServerVersion()
	if err != nil {
		return false, err
	}
	return true, nil
}

// CurrentContext returns current k8s context
func (c *ClientProvider) CurrentContext() (string, error) {
	config, err := c.config.RawConfig()
	if err != nil {
		return "", err
	}
	return config.CurrentContext, nil
}

// CurrentNamespace returns current k8s namespace
func (c *ClientProvider) CurrentNamespace() (string, error) {
	config, err := c.config.RawConfig()
	if err != nil {
		return "", err
	}
	if currentContext, ok := config.Contexts[config.CurrentContext]; ok {
		return currentContext.Namespace, nil
	}
	return "", nil
}

// Namespaces gets k8s namespaces
func (c *ClientProvider) Namespaces() ([]v1.Namespace, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.callTimout)
	defer cancel()
	namespaces, err := c.kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return namespaces.Items, nil
}

// SwitchNamespace switch to provided namespace
func (c *ClientProvider) SwitchNamespace(ns string) error {
	config, err := c.config.RawConfig()
	if err != nil {
		return err
	}
	if currentContext, ok := config.Contexts[config.CurrentContext]; ok {
		currentContext.Namespace = ns
		err := clientcmd.ModifyConfig(c.config.ConfigAccess(), config, false)
		if err != nil {
			return err
		}
	}
	return nil
}
