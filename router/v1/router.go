package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/telemetry"

	"github.com/kiichain/price-feeder/config"
	"github.com/kiichain/price-feeder/pkg/httputil"
	"github.com/kiichain/price-feeder/router/middleware"
)

// Router defines a router wrapper used for registering v1 API routes.
type Router struct {
	logger  zerolog.Logger
	cfg     config.Config
	oracle  Oracle
	metrics Metrics
}

// New creates a new router for the oracle API
func New(logger zerolog.Logger, cfg config.Config, oracle Oracle, metrics Metrics) *Router {
	return &Router{
		logger:  logger.With().Str("module", "router").Logger(),
		cfg:     cfg,
		oracle:  oracle,
		metrics: metrics,
	}
}

// RegisterRoutes register v1 API routes on the provided sub-router.
func (r *Router) RegisterRoutes(rtr *mux.Router, prefix string) {
	v1Router := rtr.PathPrefix(prefix).Subrouter()

	// build middleware chain
	mChain := middleware.Build(r.logger, r.cfg)

	// handle all preflight request
	if r.cfg.Server.EnableCORS {
		v1Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			for _, origin := range r.cfg.Server.AllowedOrigins {
				if origin == req.Header.Get("Origin") {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With",
			)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
		})
	}

	// Handle the healthz
	v1Router.Handle(
		"/healthz",
		mChain.ThenFunc(r.healthzHandler()),
	).Methods(httputil.MethodGET)

	// Handle the prices
	v1Router.Handle(
		"/prices",
		mChain.ThenFunc(r.pricesHandler()),
	).Methods(httputil.MethodGET)

	// Handle the metrics endpoint
	if r.cfg.Telemetry.Enabled {
		v1Router.Handle(
			"/metrics",
			mChain.ThenFunc(r.metricsHandler()),
		).Methods(httputil.MethodGET)
	}
}

// healthzHandler returns a handler function for the health check endpoint
func (r *Router) healthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// The response is a simple available
		resp := HealthZResponse{
			Status: StatusAvailable,
		}

		// Get the last sync time from the oracle
		resp.Oracle.LastSync = r.oracle.GetLastPriceSyncTimestamp().Format(time.RFC3339)

		// Respond on the server
		httputil.RespondWithJSON(w, http.StatusOK, resp)
	}
}

// pricesHandler returns a handler function for the prices endpoint
func (r *Router) pricesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Build a prices response
		prices := make(map[string]math.LegacyDec, len(r.oracle.GetPrices()))

		// Iterate over the prices and build the response
		for _, price := range r.oracle.GetPrices() {
			prices[price.Denom] = price.Amount
		}
		// Prepare the response
		resp := PricesResponse{
			Prices: prices,
		}

		// Respond on the server
		httputil.RespondWithJSON(w, http.StatusOK, resp)
	}
}

// metricsHandler returns a handler function for the metrics endpoint
func (r *Router) metricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get the format from the request, we use prometheus as the default
		format := strings.TrimSpace(req.FormValue("format"))

		// Set the default to prometheus if not specified
		if format == "" {
			format = telemetry.FormatPrometheus
		}

		// Gather the metrics
		gr, err := r.metrics.Gather(format)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to gather metrics: %s", err))
			return
		}

		// Respond with the gathered metrics
		w.Header().Set("Content-Type", gr.ContentType)
		_, _ = w.Write(gr.Metrics)
	}
}
