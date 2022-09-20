package traefik_plugin_cors_preflight

import (
	"context"
	"net/http"
)

type statusCodeRange struct {
	min int
	max int
}

type Config struct {
	Method string `json:"method,omitempty"`
	Code   int    `json:"code,omitempty"`
}

type CorsPreflight struct {
	name            string
	next            http.Handler
	Method          string
	Code            int
	StatusCodeRange statusCodeRange
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	_ = ctx
	return &CorsPreflight{
		next:            next,
		name:            name,
		Code:            config.Code,
		Method:          config.Method,
		StatusCodeRange: statusCodeRange{min: 100, max: 599},
	}, nil
}

func (r *CorsPreflight) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == r.Method {
		rw.WriteHeader(r.Code)
	}
}
