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
	t.Run("Unit test: test ServeHTTP force return - success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(ctx, http.MethodOptions, "http://localhost", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, 204)
	})
	t.Run("Unit test: test ServeHTTP next middleware - success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, 200)
	})
}

func TestServeFailed(t *testing.T) {
	t.Run("Unit test: test conf with method not allowed - failed", func(t *testing.T) {
		cfg := traefik_plugin_cors_preflight.CreateConfig()
		cfg.Method = "GET"
		cfg.Code = 204

		ctx := context.Background()
		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		_, err := traefik_plugin_cors_preflight.New(ctx, next, cfg, "preflight-plugin")

		assert.Equal(t, err.Error(), "method is not allowed: GET")
	})
	t.Run("Unit test: test conf with code smallest than min - failed", func(t *testing.T) {
		cfg := traefik_plugin_cors_preflight.CreateConfig()
		cfg.Method = "OPTIONS"
		cfg.Code = 99

		ctx := context.Background()
		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		_, err := traefik_plugin_cors_preflight.New(ctx, next, cfg, "preflight-plugin")

		assert.Equal(t, err.Error(), "status code is smallest than minimum value: 100")
	})
	t.Run("Unit test: test conf with code biggest than max - failed", func(t *testing.T) {
		cfg := traefik_plugin_cors_preflight.CreateConfig()
		cfg.Method = "OPTIONS"
		cfg.Code = 600

		ctx := context.Background()
		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		_, err := traefik_plugin_cors_preflight.New(ctx, next, cfg, "preflight-plugin")

		assert.Equal(t, err.Error(), "status code is biggest than maximum value: 599")
	})
}
