package main

import (
	"fmt"
	"log"

	"github.com/rohilsurana/rok/config"
	"github.com/rohilsurana/rok/kafka"
	"github.com/rohilsurana/rok/protos"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.NewConfigs("config.yaml")

	if cfg.Mode == "worker" {
		startWorker(cfg)
	} else if cfg.Mode == "server" {
		startServer(cfg)
	}
}

func startWorker(cfg *config.Config) {
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Consumer)

	for {
		consumer.Consume(func(msg *kafka.Message) error {
			m := protos.Request{}
			err := proto.Unmarshal(*msg.GetMessage(), &m)

			if err != nil {
				log.Println("Error while parsing message: ", err)
			}

			req := fasthttp.AcquireRequest()
			resp := fasthttp.AcquireResponse()
			defer fasthttp.ReleaseRequest(req)
			defer fasthttp.ReleaseResponse(resp)

			log.Println("Config service url: ", fmt.Sprintf("%s:%v", cfg.Worker.Hostname, cfg.Worker.Port))
			log.Println("Request: ", m.RequestUri)
			req.SetBody(m.Body)
			req.Header.SetMethod(m.Method)
			for _, header := range m.Headers {
				req.Header.AddBytesKV(header.Key, header.Value)
			}
			req.SetRequestURI(m.RequestUri)
			req.SetHost(fmt.Sprintf("%s:%v", cfg.Worker.Hostname, cfg.Worker.Port))
			err = fasthttp.Do(req, resp)

			if err != nil {
				log.Println("Error while requesting downstream: ", err)
			}

			if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
				log.Println("Error while requesting downstream: ", "Status: ", resp.StatusCode(), string(resp.Body()))
			}

			fmt.Println("Request sent for message - ", *msg.GetTopic(), " ", msg.GetPartition(), " ", msg.GetOffset())
			return nil
		})
	}
}

func startServer(cfg *config.Config) {
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Producer)
	defer producer.Close()

	handler := func(ctx *fasthttp.RequestCtx) {
		protoMessage := &protos.Request{
			Method:      string(ctx.Method()),
			RequestUri:  string(ctx.RequestURI()),
			HttpVersion: "HTTP/1.1",
			Headers:     []*protos.Header{},
			Body:        ctx.Request.Body(),
		}
		message, err := proto.Marshal(protoMessage)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		}

		err = producer.SendSync(message)
		if err != nil {
			log.Print("Failed to send message: ", err)
		}
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Hostname, cfg.Server.Port)
	log.Printf("Starting HTTP Server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
