package main

import (
	"context"
	conf "github.com/vllvll/diploma/internal/config"
	"github.com/vllvll/diploma/internal/handlers"
	"github.com/vllvll/diploma/internal/integrations"
	"github.com/vllvll/diploma/internal/repositories"
	"github.com/vllvll/diploma/internal/routes"
	"github.com/vllvll/diploma/internal/services"
	"github.com/vllvll/diploma/pkg/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := conf.CreateGophermartConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	db, err := postgres.ConnectDatabase(config.DatabaseURI)
	if err != nil {
		log.Fatalf("Error with database: %v", err)
	}
	defer db.Close()

	userRepository := repositories.NewUserRepository(db)
	tokenRepository := repositories.NewTokenRepository(db)
	orderRepository := repositories.NewOrderRepository(db)
	balanceRepository := repositories.NewBalanceRepository(db)

	cryptService := services.NewCrypt(config.Key)
	luhnService := services.NewLuhn()

	errCh := make(chan error)
	loyaltyCh := make(chan string)

	loyaltyClient := integrations.NewLoyaltyClient(config, loyaltyCh, errCh, orderRepository, balanceRepository)
	go loyaltyClient.Processing()

	handler := handlers.NewHandler(userRepository, tokenRepository, orderRepository, balanceRepository, cryptService, luhnService, loyaltyCh, errCh)
	router := routes.NewRouter(*handler, userRepository, tokenRepository)
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

			close(errCh)
			close(loyaltyCh)

			return
		case <-errCh:
			log.Printf("Error: %v\n", err)
		}
	}
}
