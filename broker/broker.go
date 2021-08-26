package broker

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strconv"
	"time"
)

var conn *nats.Conn

func init() {
	natsURL := os.Getenv("TP_NATS_ADDR")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.RetryOnFailedConnect(true),
	}

	var reconnectWait = 2
	if val := os.Getenv("TP_NATS_RECONNECT_WAIT"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("illegal envvar TP_NATS_RECONNECT_WAIT: %s", val)
		}
		if i > 0 {
			reconnectWait = i
		}
	}
	if reconnectWait > 0 {
		opts = append(opts, nats.ReconnectWait(time.Duration(reconnectWait)*time.Second))
	}

	if val := os.Getenv("TP_NATS_TOKEN"); val != "" {
		opts = append(opts, nats.Token(val))
	}

	username := os.Getenv("TP_NATS_USERNAME")
	password := os.Getenv("TP_NATS_PASSWORD")
	if username != "" && password != "" {
		opts = append(opts, nats.UserInfo(username, password))
	}

	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		panic(err)
	}
	conn = nc
}

func Stop() {
	if conn != nil {
		conn.Close()
	}
}

func GetConn() *nats.Conn {
	return conn
}

