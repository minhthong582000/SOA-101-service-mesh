package consul

import (
	"crypto/tls"
	"fmt"
	"net/http"

	consul "github.com/hashicorp/consul/api"
)

// Client provides an interface for getting data out of Consul
type Client interface {
	// Get a Service from consul
	Service(string, string) ([]*consul.ServiceEntry, *consul.QueryMeta, error)
	// Register a service with local agent
	Register(string, int, []string) error
	// Deregister a service with local agent
	DeRegister(string) error
}

type client struct {
	consul *consul.Client
}

// NewConsulClient returns a Client interface for given consul address
func NewConsulClient(addr string, consulTLSConfig *tls.Config) (Client, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	if consulTLSConfig != nil {
		config.Scheme = "https"
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: consulTLSConfig,
		}
	} else {
		config.Scheme = "http"
	}

	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &client{consul: c}, nil
}

// Register a service with consul local agent
func (c *client) Register(name string, port int, tags []string) error {
	reg := &consul.AgentServiceRegistration{
		ID:   name,
		Name: name,
		Port: port,
		Tags: tags,
	}

	return c.consul.Agent().ServiceRegister(reg)
}

// DeRegister a service with consul local agent
func (c *client) DeRegister(id string) error {
	return c.consul.Agent().ServiceDeregister(id)
}

// Service return a service
func (c *client) Service(service, tag string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	passingOnly := true
	addrs, meta, err := c.consul.Health().Service(service, tag, passingOnly, nil)
	if len(addrs) == 0 && err == nil {
		return nil, nil, fmt.Errorf("service ( %s ) was not found", service)
	}
	if err != nil {
		return nil, nil, err
	}

	return addrs, meta, nil
}