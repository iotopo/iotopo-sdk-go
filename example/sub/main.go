package main

import (
	"github.com/nats-io/nats.go"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/broker"
	"log"
	"os"
	"os/signal"
)

func main() {
	defer broker.Stop()
	conn := broker.GetConn()

	if sub, err := conn.Subscribe("my.data", func(msg *nats.Msg) {
		log.Println("data received:", string(msg.Data))
	}); err != nil {
		panic(err)
	} else {
		defer sub.Unsubscribe()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Subscriber ...")
}
