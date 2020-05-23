# HTTP Proxy

A small web application to receive HTTP requests and encode the whole http request into a protobuf message and push it into Kafka to be relayed over somewhere else.

### Dependencies

- fasthttp
- cleanenv
- confluent-kafka-go
