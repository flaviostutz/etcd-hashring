# etcd-hashring

[<img src="https://goreportcard.com/badge/github.com/flaviostutz/etcd-hashring"/>](https://goreportcard.com/report/github.com/flaviostutz/etcd-hashring)

Golang Lib to monitor and distribute nodes registered in ETCD dir based on a key for sharding capabilities (employ consistent hashing)

This list will connect to an ETCD server and start to watch a certain dir. When some node in this dir is updated, the internal node are updated.

When requested about the node for a arbitrary elementId, it will get the current nodes list in memory, and apply Consistentent Hash ring to select a ingle node for the element.

See a REST API example at /simple-server

## Usage - Golang lib

```golang
    import (
        hashring "github.com/flaviostutz/etcd-hashring"
    )

    //connect etcd hashring to the source etcd server
	etcdHashring, err := hashring.NewETCDHashring("http://0.0.0.0:2379", "/webservers/account", 30 * time.Second)
	if err != nil {
		return nil, err
	}

	//connect and keep nodes list updated
	go etcdHashring.Run()

    //get Node according to elementID
    node, err := GetNode(elementID)
    if err!=nil {
		return nil, err
    }

    fmt.Printf("Target node for element id %s is %s", elementID, node.Name)

```

## Usage - API Server

There is a simple API server that can be used to expose this library's results out-of-the-box.

* Create docker-compose.yml

```yml
version: '3.7'

services:

  etcd-hashring:
    image: flaviostutz/etcd-hashring
    ports:
      - "3000:3000"
    restart: always
    environment:
      - LOG_LEVEL=debug
      - ETCD_URL=http://etcd0:2379
      - ETCD_SERVICE_PATH=/webservers/test1

  etcd0:
    image: quay.io/coreos/etcd:v3.2.25
    environment:
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd0:2379

  etcd3-lucas:
    image: flaviostutz/etcd-browser-lucas:0.0.1
    ports:
      - 8888:8080
    environment:
      - ENDPOINTS=http://etcd0:2379

  sample-registrar:
    depends_on:
      - etcd0
    image: flaviostutz/etcd-registrar
    environment:
      - LOG_LEVEL=debug
      - ETCD_URL=http://etcd0:2379
      - ETCD_BASE=/webservers
      - SERVICE=test1
      - PORT=3000
      - TTL=60
      - INFO=
```

* Run ```docker-compose up -d```

* After initialization, run ```curl localhost:3000/config/abc123```

* Add more possible target instances: ```docker-compose scale sample-registrar=10```

* Run again ```curl localhost:3000/config/abc123```

* Check if returned node is different for the same element identification

* Run ```curl localhost:3000/config/xyz111```

* Run ```curl localhost:3000/config/AAAAA```

* Run ```curl localhost:3000/config/BBBBB```

* Verify that each element Id will be targeted to a different node, but the same element id will always return to the same node, unless the number of target servers is updated.

## ENVs

* LOG_LEVEL - sets the log level of verbosity. defaults to 'info'
* ETCD_URL - List of ETCD cluster URLs separated by space. Ex.: "http://etcd0:2379 http://etcd1:2379" . required
* ETCD_SERVICE_PATH - ETCD dir path to watch for nodes being added or removed. required
* ETCD_TIMEOUT_SECONDS - ETCD connection timeout. defaults to '30'

