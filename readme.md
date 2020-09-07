# JSON RabbitMQ Relay Server

Simple HTTP server that forwards JSON payload in incoming HTTP POST request to the configured AMQP server. 
Initially was created to relay Prometheus AlertManager notifications to RabbitMQ.

## Params
Configuration parameters can be passed as command line arguments or environment variables
(command line arguments take precedence)

| CMD Name | Env Variable | Description |
|----------|--------------|-------------|
| bind-addr | BIND_ADDR | HTTP server bind host:port, defaults to localhost:9712 |
| amqp-url | AMQP_URL | AMQP server URL (default: amqp://guest:guest@localhost:5672) |
| routing-key | ROUTING_KEY | AMQP Routing key (queue name) (default: test) |
 
