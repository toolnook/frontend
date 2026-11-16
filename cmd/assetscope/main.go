package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/toolnook/frontend/internal/assetscope"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8090", "Address for the AssetScope dashboard")
	flag.Parse()

	handler, err := assetscope.NewHandler()
	if err != nil {
		log.Fatalf("load dashboard assets: %v", err)
	}

	server := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("AssetScope dashboard listening on http://%s", *addr)
	log.Fatal(server.ListenAndServe())
}
