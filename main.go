package main

import (
	"flag"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel

	bindAddress string
	amqpUrl     string
	routingKey  string
)

func main() {
	log.Printf("Starting json-amqp-relay...\n")

	flag.StringVar(&bindAddress, "bind-addr", lookupEnvOrString("BIND_ADDR", "localhost:9712"),
		"Local HTTP server bind host:port")
	flag.StringVar(&amqpUrl, "amqp-url", lookupEnvOrString("AMQP_URL", "amqp://guest:guest@localhost:5672"),
		"AMQP Server URL")
	flag.StringVar(&routingKey, "routing-key", lookupEnvOrString("ROUTING_KEY", "test"),
		"AMQP Routing key (queue name)")
	flag.Parse()

	var err error

	log.Printf("Connecting to AMQP url: %s\n", amqpUrl)
	connection, err = amqp.Dial(amqpUrl)
	failOnError(err, "Unable to open AMQP connection")
	defer connection.Close()

	log.Printf("Opening AMQP channel")
	channel, err = connection.Channel()
	failOnError(err, "Unable to open AMQP channel")
	defer channel.Close()

	http.HandleFunc("/", relay)
	http.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		//returns HTTP 200 by default, good enough for simple health check
	})

	log.Printf("Starting HTTP server: %s", bindAddress)
	http.ListenAndServe(bindAddress, nil)
}

func relay(_ http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	err := channel.Publish(
		"",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish message")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
