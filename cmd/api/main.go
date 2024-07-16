package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/richardbizik/gommentary/internal/config"
	"github.com/richardbizik/gommentary/internal/database"
	"github.com/richardbizik/gommentary/internal/health"
	"github.com/richardbizik/gommentary/internal/otel"
	"github.com/richardbizik/gommentary/internal/profile"
	"github.com/richardbizik/gommentary/internal/rest/handlers"
	"github.com/richardbizik/gommentary/internal/rest/middleware"
)

var httpServer *http.Server

func main() {
	profile.InitProfile()
	conf := config.InitConfig()
	appContext, ctxCancel := context.WithCancel(context.Background())
	database, err := database.NewDatabase(conf.Database)
	if err != nil {
		panic(err)
	}

	// channel that receives signal if the application should stop
	appStop := make(chan os.Signal)

	healthCheck := health.NewHealthCheck(appContext)

	baseMux := setupMux(conf, &healthCheck, database)

	httpServer = &http.Server{
		ReadHeaderTimeout: time.Second * 5,
		Addr:              fmt.Sprintf(":%d", conf.RestApi.Port),
		Handler:           baseMux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()
	healthCheck.SetReady(true)

	handleSigterm(appStop, appContext)

	ctxCancel()

	err = httpServer.Shutdown(appContext)
	if err != nil {
		panic(err)
	}
	err = database.Close()
	if err != nil {
		panic(err)
	}
}

func handleSigterm(appStop chan os.Signal, ctx context.Context) {
	signal.Notify(appStop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-appStop
	slog.Info("Received sigterm shutting down")
}

func setupMux(conf config.Config, healthCheck *health.HealthCheck, db *database.Sqlite) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Compress(5))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Authorization(conf.JWT))
	r.Use(middleware.Cors())
	healthCheck.RegisterHandlers(r)

	if conf.EnableOTEL {
		metricsHandler, err := otel.SetupOtel(conf)
		if err != nil {
			panic(err)
		}
		r.Get("/metrics", metricsHandler.ServeHTTP)
	}
	swagger, err := handlers.GetSwagger()
	if err != nil {
		panic(err)
	}
	r.Route(config.Conf.RestApi.Context, func(r chi.Router) {
		// register metrics endpoint
		r.Use(oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(context.Context, *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		}))
		handlers.RegisterOApiHandlers(r, conf.JWT.JWTRequired, db)
	})
	return r
}
