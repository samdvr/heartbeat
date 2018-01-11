# Heartbeat

## Disclaimer

This project is a basic prototype for more battle tested service discovery libraries, checkout
[consul](https://github.com/hashicorp/consul) or [eureka](https://github.com/Netflix/eureka)

## Creating a heartbeat

Nodes should hit this endpoint every 30 seconds to be considered healthy. If a node misses 2 heartbeats (60 seconds), it will get evicted from service's healthy node list.

### POST /heartbeat

Parameters
service   | string | required | Service name of this node.
ip        | string | required | IP Address of this node.
timestamp | string | required | Epoch timestamp of the heartbeat.

Example:

#### Request

```
curl -X POST \
  http://localhost:8080/heartbeat \
  -d 'ip=127.0.0.1&service=user&timestamp=1515683663'
```

#### Response
```
  Status: 201
  {}
 ```


## Getting a healthy node for a service

### GET /service/:name
Returns a healthy node's ip address for a given service.

#### Request

```
  curl -X GET http://localhost:8080/service/sam
```

#### Response

```
  {"ip":"127.0.0.1"}
```
