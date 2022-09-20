package traefik_plugin_cors_preflight_test

import (
	"context"
	"github.com/Medzoner/traefik-plugin-cors-preflight"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeSuccess(t *testing.T) {
	cfg := traefik_plugin_cors_preflight.CreateConfig()
	cfg.Method = "OPTIONS"
	cfg.Code = 204

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin_cors_preflight.New(ctx, next, cfg, "preflight-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodOptions, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, recorder.Code, 204)
}
