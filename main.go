package main

import (
	"os"

	client "github.com/minhthong582000/SOA-101-service-mesh/client"
	server "github.com/minhthong582000/SOA-101-service-mesh/server"
)

func main() {
	if os.Getenv("TYPE") == "SERVER" {
		server.Server()
	} 
	if os.Getenv("TYPE") == "CLIENT" {
		client.Client()
	}
}