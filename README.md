# SOA-101-service-mesh

## Consul Reference Architecture

![architecture](https://learn.hashicorp.com/img/consul-arch-single.png)

## Connect-Native App Integration

Service mesh using Consul Connect-Native

[Service Mesh - Connect-Native App Integration](https://www.consul.io/docs/connect/native)

[Service Mesh - Connect Configuration](https://www.consul.io/docs/connect/configuration)

[Service Mesh - Connect-Native App Integration With Go](https://www.consul.io/docs/connect/native/go)

## DNS Interface

Consul service discovery DNS

[DNS Interface](https://www.consul.io/docs/discovery/dns)

## Demo

Run:

```bash
make run
```

Test API:

```bash
# Send get request to counting client
$ curl localhost:8081/
{"count":2,"hostname":"4b2a8b9fd305"}

# Query server address
$ curl localhost:8081/query
192.168.208.7:8080
```

![deploy](./docs/images/consul.png)

## Intentions

[Service Mesh - Intentions](https://www.consul.io/docs/connect/intentions)

Test intentions Client --x--> Server:

![intentions](./docs/images/intentions.png)

Test API:

```bash
# Send get request to counting client
$ curl localhost:8081/
Connect to server failed: Get "https://server.service.dc1.consul/count": 
remote error: tls: bad certificate
```

## Key/Value Store

Test kv store:

![kv](./docs/images/kv.png)

```bash
# Send get request to counting client
$ curl localhost:8081/kv?token=wassupbro
wassup/bro: {
  "sup": "bruh"
}

$ curl localhost:8081/kv?token=not/exist
Value not found for key: not/exist
```

