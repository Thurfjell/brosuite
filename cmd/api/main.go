package main

import (
	"brosuite/internal/api"
	"brosuite/internal/api/dashboard"
	"brosuite/internal/api/widget"
	"brosuite/internal/user"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	userService := user.New()
	widget, _ := widget.New(userService)
	dashboard, err := dashboard.New(widget, userService)

	if err != nil {
		log.Fatal(err)
	}

	server := api.Server(api.WithRoutes(dashboard.Routes), api.WithRoutes(widget.Routes))

	go func() {
		log.Printf("Starting server at %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Stopping server. Bye!")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shut down error: %v", err)
	}
}
