CoAPMQ: Publish-Subscribe Broker for the Constrained Application Protocol (CoAP) in Golang
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/coapmq/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/coapmq?status.svg)](https://godoc.org/github.com/kkdai/coapmq)  [![Build Status](https://travis-ci.org/kkdai/coapmq.svg?branch=master)](https://travis-ci.org/kkdai/coapmq)
 
    
Features
---------------

It is Golang implement base on draft RFC "[Publish-Subscribe Broker for the Constrained Application Protocol (CoAP)](https://datatracker.ietf.org/doc/draft-koster-core-coap-pubsub/?include_text=1)". It is a replace version of [CoAPMQ](https://datatracker.ietf.org/doc/draft-koster-core-coapmq/). This package base on latest draft spec (2016/01/22).


Features
---------------

- Supprt pub/sub mechanism base on CoAP
- It include a simple client/server
- Add extra heart beat mechanism to ensure UDP tunnel alive.


Install
---------------
#### Install package:
- `go get github.com/kkdai/coapmq `


#### Install binary:
- Install simple server:
	- `go get github.com/kkdai/coapmq/coapmq_server`
- Install simple interactive client: 
	- `go get github.com/kkdai/coapmq/coapmq_client`


Usage
---------------

#### Server side example

Create a 1024 buffer for pub/sub server and listen 5683 (default port for CoAP)

```go
package main

import (
	"log"

	. "github.com/kkdai/coapmq"
)

func main() {
	log.Println("Server start....")
	serv := NewBroker(1024)
	serv.ListenAndServe(":5683")
}
```

#### Client side example

Create a client to read input flag to send add/remove subscription to server.

```go
package main

import (
	"log"

	. "github.com/kkdai/coapmq"
)

func main() {
	serverAddr := "localhost:5683"
	client := NewClient(serverAddr)

	//Create Topic
	err := client.CreateTopic("topic1")
	err = client.CreateTopic("topic2")

	//Remove Topic
	err = RemoveTopic("topic2")	

	//Subsciption Topic
	ch, err := client.Subscription("topic1")
	log.Println("Wait and get sub:", <-ch)
}
```



TODO
---------------

- Hadle for UDP packet lost condition
- Gracefully network access


Benchmark
---------------
TBD

Inspired
---------------

- [CoAPMQ RFC Draft](https://datatracker.ietf.org/doc/draft-koster-core-coap-pubsub/?include_text=1)
- [RFC 7252: The Constrained Application Protocol (CoAP)](http://tools.ietf.org/html/rfc7252)
- [RFC 7641: Observing Resources in the Constrained Application Protocol (CoAP)](https://tools.ietf.org/html/rfc7641)
- [MQTT and CoAP, IoT Protocols](https://eclipse.org/community/eclipse_newsletter/2014/february/article2.php)
- [https://github.com/dustin/go-coap](https://github.com/dustin/go-coap)
- [https://github.com/gotthardp/rabbitmq-coap-pubsub](https://github.com/gotthardp/rabbitmq-coap-pubsub)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

