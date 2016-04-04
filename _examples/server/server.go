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
