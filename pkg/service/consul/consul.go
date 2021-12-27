package consul

import (
	"crypto/tls"
	"fmt"
	"net/http"

	consul "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

// Client provides an interface for getting data out of Consul
type Client interface {
	// Create connect service
	CreateConnect(string) *connect.Service
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

func (c *client) CreateConnect(name string) *connect.Service {
	// Create an instance representing this service. "name" is the
	// name of _this_ service. The service should be cleaned up via Close.
	svc, _ := connect.NewService(name, c.consul)

	return svc
}

// Register a service with consul local agent
func (c *client) Register(name string, port int, tags []string) error {
	reg := &consul.AgentServiceRegistration{
		ID:   name,
		Name: name,
		Port: port,
		Tags: tags,
	}
	reg.Connect = &consul.AgentServiceConnect{
		Native: true, // Enable this to use Connect
	}

	// Register Health check
	// reg.Check = & consul.AgentServiceCheck { 
	// 	HTTP: fmt.Sprintf("https://%s:%d%s", name, port, "/ping"),
	// 	Timeout: "3s", // Timeout
	// 	Interval: "5s", // Health check interval
	// 	TLSSkipVerify: true, // Skip TLS Verify
	// 	// After 30 seconds after failure, 
	// 	// delete this service, logout, equivalent to expiration time
	// 	DeregisterCriticalServiceAfter: "5m", 
	// 	// GRPC support, execute the address of the health check, 
	// 	// Service will be transmitted to the Health.Check function
	// 	// Grpc: fmt.sprintf ("% v:% v /% v", IP, r.Port, r.service),
	// }

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