package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/consul/connect"
	consul "github.com/minhthong582000/SOA-101-service-mesh/pkg/service/consul"
	kv "github.com/minhthong582000/SOA-101-service-mesh/pkg/service/kv"
)

type ClientHandler struct {
	svc *connect.Service
	kvService kv.Client
	uServiceConsul consul.Client
}

// Counting send request to counting server (counting service)
func (c ClientHandler) Counting(rw http.ResponseWriter, r *http.Request) {
	// Get an HTTP client
	httpClient := c.svc.HTTPClient()

	// Perform a request to server, then use the standard response
	serverEnpoint := "https://" + os.Getenv("CONSUL_QUERY_APP_ID") + ".service.dc1.consul/count"
	resp, err := httpClient.Get(serverEnpoint)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte("Connect to server failed: " + err.Error()))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte("Error reading response from counting server !"))
		return
	}

	rw.WriteHeader(200)
	rw.Write(body)
}

func (c ClientHandler) GetKVStore(rw http.ResponseWriter, r *http.Request) {
	k := r.URL.Query().Get("token")

	v, err := c.kvService.Get(k)
	if err != nil || v == "" {
		rw.WriteHeader(404)
		rw.Write([]byte("Value not found for key: " + k))
		return
	}

	rw.WriteHeader(200)
	rw.Write([]byte(k + ": " + v))
}

func (c ClientHandler) QueryConsul(rw http.ResponseWriter, r *http.Request) {
	services, _, err := c.uServiceConsul.Service(os.Getenv("CONSUL_QUERY_APP_ID"), "demo")
	if err != nil {
		fmt.Println("Discover failed:", err)
		rw.WriteHeader(500)
		rw.Write([]byte("Discover failed!"))
		return
	}

	log.Println("Found service at these locations:")
	var adrr []string 
	for _, v := range services {
		log.Println(fmt.Sprintf("%s:%d", v.Node.Address, v.Service.Port))
		adrr = append(adrr, fmt.Sprintf("%s:%d", v.Node.Address, v.Service.Port))
	}
	
	rw.WriteHeader(200)
	rw.Write([]byte(strings.Join(adrr," ")))
}

// heathCheck handle health check
func (c ClientHandler) heathCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Client create a counting client app
func Client() {
	consulAdrr := os.Getenv("CONSUL_AGENT_ADDR")

	// Uncomment to enable tls
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

	// New consul client
	uServiceConsul, err := consul.NewConsulClient(consulAdrr, nil)
	if err != nil {
		log.Fatalln("Can't find consul:", err)
	}

	// New consul kv client
	kvService, err := kv.NewKVClient(consulAdrr, "/", nil)
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

	// Create connect http
	svc := uServiceConsul.CreateConnect(svcName)
	defer svc.Close()

	clientHandler := &ClientHandler{
		svc: svc, 
		kvService: kvService,
		uServiceConsul: uServiceConsul,
	}
	// Send request to counting server (counting service)
	http.HandleFunc("/", clientHandler.Counting)
	// Handle health check
	http.HandleFunc("/ping", clientHandler.heathCheck)
	// Get data from consul kv
	http.HandleFunc("/kv", clientHandler.GetKVStore)
	// Query server information
	http.HandleFunc("/query", clientHandler.QueryConsul)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}