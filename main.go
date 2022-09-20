package traefik_plugin_cors_preflight

import (
	"context"
	"fmt"
	"net/http"
)

type statusCodeRange struct {
	min int
	max int
}

type Config struct {
	Method          string `json:"method,omitempty"`
	Code            int    `json:"code,omitempty"`
	StatusCodeRange statusCodeRange
	AllowMethods    []string
}

type CorsPreflight struct {
	name   string
	next   http.Handler
	Method string
	Code   int
}

func CreateConfig() *Config {
	return &Config{
		StatusCodeRange: statusCodeRange{min: 100, max: 599},
		AllowMethods:    []string{"OPTIONS"},
	}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	_ = ctx
	if config.Code < config.StatusCodeRange.min {
		return nil, fmt.Errorf("status code is smallest than minimum value: %v", config.StatusCodeRange.min)
	}
	if config.Code > config.StatusCodeRange.max {
		return nil, fmt.Errorf("status code is biggest than maximum value: %v", config.StatusCodeRange.max)
	}
	if !contains(config.AllowMethods, config.Method) {
		return nil, fmt.Errorf("method is not allowed: " + config.Method)
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
