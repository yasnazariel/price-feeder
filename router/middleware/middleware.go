package middleware

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/kiichain/price-feeder/config"
)

// Build constructs a middleware chain for the HTTP server based on the provided
func Build(logger zerolog.Logger, cfg config.Config) alice.Chain {
	// Build the middleware chain
	mChain := alice.New()
	mChain = AddRequestLoggingMiddleware(mChain, logger)

	// Check if cors should be added
	if cfg.Server.EnableCORS {
		mChain = AddCORSMiddleware(mChain, logger, cfg)
	}

	return mChain
}

// AddRequestLoggingMiddleware appends HTTP logging middleware to a provided
// middleware chain.
func AddRequestLoggingMiddleware(mChain alice.Chain, logger zerolog.Logger) alice.Chain {
	mChain = mChain.Append(hlog.NewHandler(logger))
	mChain = mChain.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	mChain = mChain.Append(hlog.RequestHandler("req"))
	mChain = mChain.Append(hlog.RemoteAddrHandler("ip"))
	mChain = mChain.Append(hlog.UserAgentHandler("ua"))
	mChain = mChain.Append(hlog.RefererHandler("ref"))
	mChain = mChain.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return mChain
}

// AddCORSMiddleware appends CORS middleware to a provided middleware chain.
func AddCORSMiddleware(mChain alice.Chain, logger zerolog.Logger, cfg config.Config) alice.Chain {
	// Cors options
	opts := cors.Options{
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodOptions,
		},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Content-Type",
			"Access-Control-Allow-Headers",
			"Authorization",
			"X-Requested-With",
		},
		AllowedOrigins: cfg.Server.AllowedOrigins,
	}

	c := cors.New(opts)
	c.Log = &logger

	mChain = mChain.Append(c.Handler)

	return mChain
}
