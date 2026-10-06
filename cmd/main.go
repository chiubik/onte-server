package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"x/onte-server/internal"
	"x/onte-server/internal/handler"

	"github.com/mdp/qrterminal"
)

func main() {
	ch := make(chan string, 1)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go internal.CloudtunnelRun(ctx, ch)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.IdentifierHandler)
	mux.HandleFunc("POST /upload", handler.UploadHandler)
	mux.HandleFunc("GET /donwload", handler.DownloadHandler)

	srv := &http.Server{
		Addr:    ":3333",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server is running...")
		url := <-ch
		if url == "" {
			log.Fatal("Didn't get url")
		}
		identifier := internal.GenerateIdentifier()
		fmt.Println("Url: " + url)
		fmt.Println("Identifier: " + identifier)
		data := url + "+" + identifier
		qrterminal.Generate(data, qrterminal.L, os.Stdout)
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println("error: %s", err)
		}
	}()

	<-ctx.Done()
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("error: %s", err)
	}

}
