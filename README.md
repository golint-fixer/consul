# consul [![Build Status](https://travis-ci.org/vinxi/consul.png)](https://travis-ci.org/vinxi/consul) [![GoDoc](https://godoc.org/github.com/vinxi/consul?status.svg)](https://godoc.org/github.com/vinxi/consul) [![Coverage Status](https://coveralls.io/repos/github/vinxi/consul/badge.svg?branch=master)](https://coveralls.io/github/vinxi/consul?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/consul)](https://goreportcard.com/report/github.com/vinxi/consul) [![API](https://img.shields.io/badge/vinxi-core-green.svg?style=flat)](https://godoc.org/github.com/vinxi/consul) 

[Consul](https://www.consul.io) plugin for simple dynamic service discovery and optional traffic balancing in vinxi proxies.

Supports multiple Consul servers, non-blocking multi-thread service discovery, retry policy on discovery error and transparent fallback to the next Consul server available.

## Installation

```bash
go get -u gopkg.in/vinxi/consul.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/consul) reference.

## Example

#### Server discovery with default roundrobin balancer

```go
package main

import (
  "fmt"
  "gopkg.in/vinxi/consul.v0"
  "gopkg.in/vinxi/forward.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create the Consul client
  cc := consul.New(consul.NewConfig("web", "http://demo.consul.io"))

  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})
  vs.Use(cc)

  fw, _ := forward.New(forward.PassHostHeader(true))
  vs.UseFinalHandler(fw)

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

## License

MIT
