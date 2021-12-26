package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	consul "github.com/minhthong582000/SOA-101-service-mesh/pkg/service/consul"
)

// Count stores a number that is being counted and other data to
// return as JSON in the API.
type Count struct {
	Count    uint64 `json:"count"`
	Hostname string `json:"hostname"`
}

// CountHandler serves a JSON feed that contains a number that increments each time
// the API is called.
type CountHandler struct {
	index *uint64
}

// counter Handle counting request from client service
func (h CountHandler) counter(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(h.index, 1)
	hostname, _ := os.Hostname()
	index := atomic.LoadUint64(h.index)

	count := Count{Count: index, Hostname: hostname}

	responseJSON, _ := json.Marshal(count)
	w.WriteHeader(200)
	w.Write([]byte(responseJSON))
}

// heathCheck handle health check
func heathCheck(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// Server create a counting server
func Server() {
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

	// Create connect http
	svc := uServiceConsul.CreateConnect(svcName)
	defer svc.Close()

	//Create the default mux
	mux := http.NewServeMux()

	var index uint64
	mux.HandleFunc("/count", CountHandler{index: &index}.counter)
	mux.HandleFunc("/ping", heathCheck) // Handle health check

	// Creating an HTTP server that serves via Connect
	server := &http.Server{
		Addr:      ":8080",
		TLSConfig: svc.ServerTLSConfig(),
		Handler: mux,
	}
	// Serve!
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}