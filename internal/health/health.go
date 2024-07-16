package health

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/go-chi/chi/v5"
)

type HealthCheckFunc func() error

type HealthCheck struct {
	Funcs      []HealthCheckFunc
	ready      bool
	readyMutex *sync.Mutex
	ctx        context.Context
}

func NewHealthCheck(ctx context.Context, healthCheckFuncs ...HealthCheckFunc) HealthCheck {
	return HealthCheck{
		Funcs:      healthCheckFuncs,
		ready:      false,
		readyMutex: &sync.Mutex{},
		ctx:        ctx,
	}
}

func (hc *HealthCheck) RegisterHandlers(r chi.Router) {
	r.Get("/readyz", hc.ReadinessHandler)
	r.Get("/healthz", hc.HealthHandler)
}

func (hc *HealthCheck) HealthHandler(w http.ResponseWriter, r *http.Request) {
	var accErr error
	for _, healthCheckFunc := range hc.Funcs {
		err := healthCheckFunc()
		if err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(healthCheckFunc).Pointer()).Name()
			accErr = errors.Join(accErr, fmt.Errorf("%s: %w", funcName, err))
		}
	}
	if accErr == nil {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(503)
	w.Write([]byte(accErr.Error()))
}

func (hc *HealthCheck) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if hc.ready && hc.ctx.Err() == nil {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(503)
}

func (hc *HealthCheck) SetReady(ready bool) {
	hc.readyMutex.Lock()
	defer hc.readyMutex.Unlock()
	hc.ready = ready
}
