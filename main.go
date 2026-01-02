package traefik_plugin_cors_preflight

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
)

type statusCodeRange struct {
	Min int
	Max int
}

type Config struct {
	Method          string `json:"method,omitempty"`
	AllowMethods    []string
	StatusCodeRange statusCodeRange
	Code            int  `json:"code,omitempty"`
	Debug           bool `json:"debug,omitempty"`
}

type CorsPreflight struct {
	name   string
	next   http.Handler
	Method string
	Code   int
	Debug  bool
}

func CreateConfig() *Config {
	return &Config{
		StatusCodeRange: statusCodeRange{Min: 100, Max: 599},
		AllowMethods:    []string{http.MethodOptions},
		Method:          http.MethodOptions,
		Code:            http.StatusNoContent,
		Debug:           false,
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

	log.Printf("Plugin traefik-plugin-cors-preflight - Init with return code %d for method %s\n", config.Code, config.Method)

	return &CorsPreflight{
		next:   next,
		name:   name,
		Code:   config.Code,
		Method: config.Method,
		Debug:  config.Debug,
	}, nil
}

func (r *CorsPreflight) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.Debug {
		log.Printf("Plugin traefik-plugin-cors-preflight - Received request with method: %s\n", req.Method)
	}

	if req.Method == r.Method {
		if r.Debug {
			log.Printf("Plugin traefik-plugin-cors-preflight - Returning status code: %d for method: %s\n", r.Code, r.Method)
		}
		rw.WriteHeader(r.Code)
		return
	}

	if req.Method != r.Method {
		if r.Debug {
			log.Printf("Plugin traefik-plugin-cors-preflight - Passing to next middleware for method: %s\n", req.Method)
		}
		r.next.ServeHTTP(rw, req)
	}
}
