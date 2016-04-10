package main

import (
	"fmt"
	"gopkg.in/vinxi/apachelog.v0"
	"gopkg.in/vinxi/consul.v0"
	"gopkg.in/vinxi/forward.v0"
	"gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
	// Create a new vinxi proxy
	vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})
	vs.Use(apachelog.Default)

	// Create and attach Consul client
	vs.Use(consul.New(consul.NewConfig("web", "http://demo.consul.io")))

	fw, _ := forward.New(forward.PassHostHeader(true))
	vs.UseFinalHandler(fw)

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
