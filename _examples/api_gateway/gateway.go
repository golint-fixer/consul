package main

import (
	"fmt"
	"gopkg.in/vinxi/consul.v0"
	"gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
	// Create the Consul client for web service
	web := consul.New(consul.NewConfig("web", "http://demo.consul.io"))

	// Create the Consul client for proxy service
	consul := consul.New(consul.NewConfig("consul", "http://demo.consul.io"))

	// Create a new vinxi proxy
	vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

	vs.Vinxi.Get("/").Use(web)
	vs.Vinxi.Get("/consul").Use(consul)

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
