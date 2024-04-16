package main

import (
	"go-web-backend/handle"
	"go-web-backend/servercontext"
	"log"
	"net/http"
)

func main() {
	srvctx := servercontext.New()
	http.HandleFunc("GET /hello", srvctx.Wrap(handle.Hello))
	log.Fatalln(http.ListenAndServe(":8181", nil))
}
