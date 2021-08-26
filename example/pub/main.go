package main

import "gogs.iotopo.com/iotopo/iotopo-sdk-go/broker"

func main() {
	defer broker.Stop()
	conn := broker.GetConn()
	if err := conn.Publish("my.subject", []byte("hello world")); err != nil {
		panic(err)
	}
}
