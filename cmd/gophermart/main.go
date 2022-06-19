package main

import (
	"context"
	conf "github.com/vllvll/diploma/internal/config"
	"github.com/vllvll/diploma/internal/handlers"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/routes"
	"github.com/vllvll/diploma/pkg/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := conf.CreateServerConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	db, err := postgres.ConnectDatabase(config.DatabaseUri)
	if err != nil {
		log.Fatalf("Error with database: %v", err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)

	handler := handlers.NewHandler(userRepository)
	router := routes.NewRouter(*handler)
	//router = routes.NewRouter(*handler)
	router.RegisterHandlers()

	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: router.Router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error with HTTP server ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			log.Println("Graceful shutdown")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if err := httpServer.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			cancel()

			return
		}
	}
}
