package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"x/onte-server/internal"
	"x/onte-server/internal/handler"

	"github.com/mdp/qrterminal"
)

func main() {
	identifier := internal.GenerateIdentifier()
	ch := make(chan string, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go internal.CloudtunnelRun(ctx, ch)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.IdentifierHandler)
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		handler.UploadHandler(w, r, identifier)
	})
	mux.HandleFunc("GET /donwload", func(w http.ResponseWriter, r *http.Request) {
		handler.DownloadHandler(w, r, identifier)
	})

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
	wait, cancel := context.WithTimeout(context.Background(), 30*time.Second) //added a timeout if there will be a problem to shutdown the server
	defer cancel()
	err := srv.Shutdown(wait)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
