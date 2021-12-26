package kv

import (
	"crypto/tls"
	"net/http"

	consul "github.com/hashicorp/consul/api"
)

// Client provides an interface for getting data out of Consul
type Client interface {
	// Create connect service
	Set(k string, v string) error
	// Get a Service from consul
	Get(k string) (v string, err error)
	// Register a service with local agent
	Delete(k string) error
}

// client is a gokv.Store implementation for Consul.
type client struct {
	c      *consul.KV
	folder string
}

// Options are the options for the Consul client.
type Options struct {
	// URI scheme for the Consul server.
	// Optional ("http" by default).
	Scheme string
	// Address of the Consul server, including port number.
	// Optional ("127.0.0.1:8500" by default).
	Address string
	// Directory under which to store the key-value pairs.
	// The Consul UI calls this "folder".
	// Optional (none by default).
	Folder string
}

// NewClient creates a new Consul client.
func NewKVClient(addr string, folder string, consulTLSConfig *tls.Config) (Client, error) {
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

	return &client{c: c.KV(), folder: folder}, nil
}

// Set stores the given value for the given key.
func (c client) Set(k string, v string) error {
	if c.folder != "" {
		k = c.folder + "/" + k
	}
	kvPair := consul.KVPair{
		Key:   k,
		Value: []byte(v),
	}
	_, err := c.c.Put(&kvPair, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves the stored value for the given key.
func (c client) Get(k string) (v string, err error) {

	if c.folder != "" {
		k = c.folder + "/" + k
	}
	kvPair, _, err := c.c.Get(k, nil)
	if err != nil {
		return "", err
	}
	// If no value was found return empty
	if kvPair == nil {
		return "", nil
	}
	data := kvPair.Value

	return string(data), nil
}

// Delete deletes the stored value for the given key.
func (c client) Delete(k string) error {
	if c.folder != "" {
		k = c.folder + "/" + k
	}
	_, err := c.c.Delete(k, nil)

	return err
}