package main

import (
	"fmt"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/broker"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/example/rpc/helloworld"
	"log"
)

func main() {
	defer broker.Stop()
	nc := broker.GetConn()

	// This is our generated client.
	cli := helloworld.NewGreeterClient(nc)

	// Contact the server and print out its response.
	resp, err := cli.SayHello(helloworld.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Printf("Greeting: %s\n", resp.GetMessage())
}
