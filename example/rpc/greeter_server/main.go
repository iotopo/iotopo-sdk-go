package main

import (
	"context"
	"fmt"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/broker"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/example/rpc/helloworld"

	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"log"
	"os"
	"os/signal"
)

// server implements the helloworld.GreeterServer interface.
type server struct{}

// SayHello is an implementation of the SayHello method from the definition of
// the Greeter service.
func (s *server) SayHello(ctx context.Context, req helloworld.HelloRequest) (resp helloworld.HelloReply, err error) {
	resp.Message = "Hello " + req.Name
	return
}

func main() {
	defer broker.Stop()
	nc := broker.GetConn()

	// Our server implementation.
	s := &server{}

	// The NATS handler from the helloworld.nrpc.proto file.
	h := helloworld.NewGreeterHandler(context.TODO(), nc, s)

	// Start a NATS subscription using the handler. You can also use the
	// QueueSubscribe() method for a load-balanced set of servers.
	//sub, err := nc.Subscribe(h.Subject(), h.Handler)
	sub, err := nc.QueueSubscribe(h.Subject(), "greeter", h.Handler)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Keep running until ^C.
	fmt.Println("server is running, ^C quits.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(c)
}
