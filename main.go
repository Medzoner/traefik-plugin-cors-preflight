package traefik_plugin_cors_preflight

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/rs/zerolog"
)

type headerName string

const (
	headerAllowOrigin  headerName = "Access-Control-Allow-Origin"
	headerAllowMethods headerName = "Access-Control-Allow-Methods"
	headerAllowHeaders headerName = "Access-Control-Allow-Headers"
)

func (h headerName) String() string {
	return string(h)
}

type (
	statusCodeRange struct {
		Min int
		Max int
	}

	Config struct {
		Method          string   `json:"method,omitempty"`
		AllowOrigins    []string `json:"allowOrigins,omitempty"`
		AllowMethods    []string `json:"allowMethods,omitempty"`
		AllowHeaders    []string `json:"allowHeaders,omitempty"`
		StatusCodeRange statusCodeRange
		Code            int  `json:"code,omitempty"`
		Debug           bool `json:"debug,omitempty"`
	}

	CorsPreflight struct {
		next         http.Handler
		name         string
		Method       string
		AllowOrigins []string
		AllowMethods []string
		AllowHeaders []string
		Code         int
	}
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(zerolog.ErrorLevel)

func CreateConfig() *Config {
	return &Config{
		StatusCodeRange: statusCodeRange{Min: 100, Max: 599},
		AllowMethods:    []string{http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowHeaders:    []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
		Method:          http.MethodOptions,
		Code:            http.StatusNoContent,
		Debug:           false,
		AllowOrigins:    []string{"*"},
	}
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Code < config.StatusCodeRange.Min {
		return nil, fmt.Errorf("status code is smallest than minimum value: %d", config.StatusCodeRange.Min)
	}
	if config.Code > config.StatusCodeRange.Max {
		return nil, fmt.Errorf("status code is biggest than maximum value: %d", config.StatusCodeRange.Max)
	}
	if !slices.Contains(config.AllowMethods, config.Method) {
		return nil, fmt.Errorf("method is not allowed: %s", config.Method)
	}

	if config.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}

	log("Plugin traefik-plugin-cors-preflight - Init with return code %d for method %s", config.Code, config.Method)

	return &CorsPreflight{
		next:         next,
		name:         name,
		Code:         config.Code,
		Method:       config.Method,
		AllowOrigins: config.AllowOrigins,
		AllowMethods: config.AllowMethods,
		AllowHeaders: config.AllowHeaders,
	}, nil
}

func (r *CorsPreflight) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log("Plugin traefik-plugin-cors-preflight - Received request with method: %s", req.Method)

	origin := req.Header.Get("Origin")
	for _, allowed := range r.AllowOrigins {
		if origin == allowed || allowed == "*" {
			rw.Header().Set(headerAllowOrigin.String(), allowed)
			break
		}
	}

	if len(r.AllowMethods) > 0 {
		rw.Header().Set(headerAllowMethods.String(), fmt.Sprintf("%s", r.AllowMethods))
	}

	if len(r.AllowHeaders) > 0 {
		rw.Header().Set(headerAllowHeaders.String(), fmt.Sprintf("%s", r.AllowHeaders))
	}

	if req.Method == r.Method {
		log("Plugin traefik-plugin-cors-preflight - Returning status code: %d for method: %s", r.Code, r.Method)
		rw.WriteHeader(r.Code)

		return
	}

	log("Plugin traefik-plugin-cors-preflight - Passing to next middleware for method: %s", req.Method)
	r.next.ServeHTTP(rw, req)
}

func log(s string, args ...any) {
	logger.Debug().Msgf(s, args...)
}
