package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"x/flap/handler"
	"x/flap/pkg"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go pkg.CloudtunnelRun(ctx)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.IdentifierHandler)
	mux.HandleFunc("POST /upload", handler.UploadHandler)
	mux.HandleFunc("GET /donwload", handler.DownloadHandler)

	fmt.Print("Server is running...")
	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
