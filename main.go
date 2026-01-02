package traefik_plugin_cors_preflight

import (
	"context"
	"fmt"
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
	Code            int `json:"code,omitempty"`
}

type CorsPreflight struct {
	name   string
	next   http.Handler
	Method string
	Code   int
}

func CreateConfig() *Config {
	return &Config{
		StatusCodeRange: statusCodeRange{Min: 100, Max: 599},
		AllowMethods:    []string{http.MethodOptions},
	}
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Code < config.StatusCodeRange.Min {
		return nil, fmt.Errorf("status code is smallest than minimum value: %v", config.StatusCodeRange.Min)
	}
	if config.Code > config.StatusCodeRange.Max {
		return nil, fmt.Errorf("status code is biggest than maximum value: %v", config.StatusCodeRange.Max)
	}
	if !slices.Contains(config.AllowMethods, config.Method) {
		return nil, fmt.Errorf("method is not allowed: %v", config.Method)
	}

	fmt.Printf("Plugin traefik-plugin-cors-preflight - Init with return code %v for method %v\n", config.Method, config.Code)

	return &CorsPreflight{
		next:   next,
		name:   name,
		Code:   config.Code,
		Method: config.Method,
	}, nil
}

func (r *CorsPreflight) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == r.Method {
		rw.WriteHeader(r.Code)
		return
	}

	if req.Method != r.Method {
		r.next.ServeHTTP(rw, req)
	}
}
