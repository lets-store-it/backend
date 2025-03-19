package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evevseev/storeit/backend/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	port = flag.String("port", "8080", "Port to listen on")
)

func main() {
	flag.Parse()

	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handler := &Handler{}
	server, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	e.Any("/*", echo.WrapHandler(server))

	addr := fmt.Sprintf(":%s", *port)
	go func() {
		if err := e.Start(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

type Handler struct {
	api.UnimplementedHandler
}

func (h *Handler) GetHealth(ctx context.Context) (api.GetHealthRes, error) {
	return &api.GetHealthOK{
		Status: api.NewOptGetHealthOKStatus(api.GetHealthOKStatusOK),
	}, nil
}
