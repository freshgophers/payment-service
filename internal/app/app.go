package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"payment-service/internal/service/catalogue"
	"payment-service/internal/service/payment"
	"payment-service/pkg/epay"
	"syscall"
	"time"

	"go.uber.org/zap"

	"payment-service/internal/config"
	"payment-service/internal/handler"
	"payment-service/internal/repository"
	"payment-service/pkg/log"
	"payment-service/pkg/server"
)

const (
	schema      = "payment"
	version     = "1.0.0"
	description = "payment-service"
	configsDir  = "configs"
)

// Run initializes whole application.
func Run() {
	logger := log.New(version, description)

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIG", zap.Error(err))
		return
	}

	repositories, err := repository.New(
		//repository.WithPostgresStore(schema, configs.POSTGRES.DSN))
		repository.WithMemoryStore())
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORY", zap.Error(err))
		return
	}
	defer repositories.Close()

	catalogueService, err := catalogue.New(
		catalogue.WithCategoryRepository(repositories.Category),
		catalogue.WithProductRepository(repositories.Product),
		catalogue.WithCategoryCache(repositories.Category),
		catalogue.WithProductCache(repositories.Product),
	)

	if err != nil {
		logger.Error("ERR_INIT_CATALOGUE_SERVICE", zap.Error(err))
		return
	}

	paymentService, err := payment.New(
		payment.WithBillingRepository(repositories.Billing),
		payment.WithBillingCache(repositories.Billing),
	)

	if err != nil {
		logger.Error("ERR_INIT_PAYMENT_SERVICE", zap.Error(err))
		return
	}

	ePayClient := epay.NewClient(epay.Credential{
		TerminalID:    configs.EPay.TerminalID,
		ClientID:      configs.EPay.ClientID,
		ClientSecret:  configs.EPay.ClientSecret,
		OauthEndpoint: configs.EPay.OauthEndpoint,
		Endpoint:      configs.EPay.Endpoint,
		JSLink:        configs.EPay.JSLink,
		BackLink:      configs.EPay.BackLink,
		PostLink:      configs.EPay.PostLink,
		Amount:        configs.EPay.Amount,
	})

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:          configs,
			CatalogueService: catalogueService,
			EPayClient:       ePayClient,
			PaymentService:   paymentService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLER", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.HTTP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVER", zap.Error(err))
		return
	}

	// Run our server in a goroutine so that it doesn't block.
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVER", zap.Error(err))
		return
	}

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("Server was successful shutdown.")
}
