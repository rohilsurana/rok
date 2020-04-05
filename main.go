package main

import (
	"fmt"
	"github.com/rohilsurana/http-proxy/config"
	"github.com/valyala/fasthttp"
	"log"
)

func main() {
	cfg := config.Configs("config.yaml")

	handler := func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())
		path := string(ctx.Path())
		headers := ctx.Request.Header.RawHeaders()
		body := ctx.PostBody()

		fmt.Println("----------------REQUEST START----------------")
		fmt.Printf("%s %s HTTP/1.1\n", method, path)
		fmt.Printf(string(headers))
		fmt.Printf(string(body))
		fmt.Println("-----------------REQUEST END-----------------")

		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting HTTP Server on %s", addr)
	if err := fasthttp.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
