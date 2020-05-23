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

}

func startServer(cfg *config.Config) {
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Producer)

	handler := func(ctx *fasthttp.RequestCtx) {
		// request := ctx.Request.String()
		protoMessage := &protos.Request{
			Method:      string(ctx.Method()),
			Path:        string(ctx.Path()),
			HttpVersion: "HTTP/1.1",
			Headers:     []*protos.Header{},
			Body:        ctx.Request.Body(),
		}
		message, err := proto.Marshal(protoMessage)
		if err != nil {
			log.Fatalln("Failed to encode address book:", err)
		}
		// fmt.Println(request)
		err = producer.SendSync(message)
		if err != nil {
			log.Print("Failed to send message: ", err)
		}
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting HTTP Server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
