CoAPMQ:
==================

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/coapmq/master/LICENSE)  [![GoDoc](https://godoc.org/github.com/kkdai/coapmq?status.svg)](https://godoc.org/github.com/kkdai/coapmq)  [![Build Status](https://travis-ci.org/kkdai/coapmq.svg?branch=master)](https://travis-ci.org/kkdai/coapmq)
    





Note
---------------

It will keep a heart beat signal from client to server if you subscription a topic to remain your UDP port channel.

Install
---------------
`go get github.com/kkdai/coapmq `


Usage
---------------

#### Server side example

Create a 1024 buffer for pub/sub server and listen 5683 (default port for CoAP)

```go

```

#### Client side example

Create a client to read input flag to send add/remove subscription to server.

```go
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
- [MQTT and CoAP, IoT Protocols](https://eclipse.org/community/eclipse_newsletter/2014/february/article2.php)
- [RFC 7252](http://tools.ietf.org/html/rfc7252)
- [https://github.com/dustin/go-coap](https://github.com/dustin/go-coap)

Project52
---------------

It is one of my [project 52](https://github.com/kkdai/project52).


License
---------------

This package is licensed under MIT license. See LICENSE for details.

