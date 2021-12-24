package main

import (
	"fmt"
	"log"
	"os"
	"time"

	consul "github.com/minhthong582000/SOA-101-service-mesh/pkg/service/consul"
)

func main() {
	consulAdrr := os.Getenv("CONSUL_AGENT_ADDR")

	// consulTLSConfig, err := consul.SetupTLSConfig(&consul.TLSConfig{
	// 	Address:            consulAdrr,
	// 	CAFile:             "/etc/tls/consul/ca.pem",
	// 	CertFile:           "/etc/tls/consul/consul.pem",
	// 	KeyFile:            "/etc/tls/consul/consul-key.pem",
	// 	// InsecureSkipVerify: true,
	// }) 
	// if err != nil {
	// 	log.Fatalln(fmt.Sprintf("SSL Configuration error: %s\n", err))
	// }

	uServiceConsul, err := consul.NewConsulClient(consulAdrr, nil)
	if err != nil {
		log.Fatalln("Can't find consul:", err)
	}

	// Register itself
	svcName := os.Getenv("CONSUL_APP_ID")
	svcPort := 8080
	svcTags := []string{"demo"}
	err = uServiceConsul.Register(svcName, svcPort, svcTags)
	if err != nil {
		log.Fatalln("Register failed:", err)
	}
	defer func() {
		time.Sleep(60 * time.Second)
		err = uServiceConsul.DeRegister(svcName)
		if err != nil {
			log.Fatalln("DeRegister failed:", err)
		}
	}()

	// Query health information
	services, _, err := uServiceConsul.Service(os.Getenv("CONSUL_QUERY_APP_ID"), "demo")
	if err != nil {
		log.Fatalln("Discover failed:", err)
	}

	log.Println("Found service at these locations:")
	for _, v := range services {
		log.Println(fmt.Sprintf("%s:%d", v.Node.Address, v.Service.Port))
	}
}