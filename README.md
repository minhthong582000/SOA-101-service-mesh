# SOA-101-service-mesh

## Consul Reference Architecture

![architecture](https://learn.hashicorp.com/img/consul-arch-single.png)

## Connect-Native App Integration

Service mesh using Consul Connect-Native

[Connect-Native App Integration](https://www.consul.io/docs/connect/native)

[Connect Configuration](https://www.consul.io/docs/connect/configuration)

[Connect-Native App Integration With Go](https://www.consul.io/docs/connect/native/go)

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
# Send get request to counting server
$ curl localhost:8081/
{"count":2,"hostname":"4b2a8b9fd305"}

# Query server address
$ curl localhost:8081/query
192.168.208.7:8080
```

![deploy](./docs/images/consul.png)
