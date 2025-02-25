package e2e

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/richardbizik/gommentary/internal/config"
	"github.com/richardbizik/gommentary/internal/database"
	"github.com/richardbizik/gommentary/internal/otel"
	"github.com/richardbizik/gommentary/internal/profile"
	"github.com/richardbizik/gommentary/internal/rest/handlers"
	"github.com/richardbizik/gommentary/internal/rest/middleware"
)

type Application struct {
	port    int
	cancel  context.CancelFunc
	appStop chan os.Signal
	server  *http.Server
}

func NewApplication(port int) *Application {
	app := Application{}
	profile.Current = profile.TEST
	os.Setenv("CONFIG_FILE", path.Join("..", "..", "conf", "api", "conf-test.yaml"))
	conf := config.InitConfig()
	conf.Database.File = path.Join(os.TempDir(), fmt.Sprintf("gommentary-%d.db", port))
	conf.Database.MigrationsDir = path.Join("..", "..", "sql", "migrations")
	_, app.cancel = context.WithCancel(context.Background())
	database, err := database.NewDatabase(conf.Database)
	if err != nil {
		panic(err)
	}
	app.appStop = make(chan os.Signal)

	baseMux := setupMux(conf, database)

	app.server = &http.Server{
		ReadHeaderTimeout: time.Second * 5,
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           baseMux,
	}
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(err.Error())
		}
	}()
	return &app
}

func (a *Application) Stop() {
	a.cancel()
	a.server.Close()
}

func setupMux(conf config.Config, db *database.Sqlite) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.Compress(5))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Authorization(conf.JWT))
	r.Use(middleware.Cors())

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

func getRandomPort() int {
	min := 50000
	max := 65000
	return rand.Intn(max-min) + min
}
