package traefik_plugin_cors_preflight_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	traefik_plugin_cors_preflight "github.com/Medzoner/traefik-plugin-cors-preflight"
	"gotest.tools/assert"
)

type testCase struct {
	name     string
	method   string
	in       in
	err      error
	expected int
}

type in struct {
	method string
}

func TestServe(t *testing.T) {
	cfg := traefik_plugin_cors_preflight.CreateConfig()
	for _, tc := range []testCase{
		{
			name:   "force return - success",
			method: http.MethodOptions,
			in: in{
				method: http.MethodOptions,
			},
			err:      nil,
			expected: http.StatusNoContent,
		},
		{
			name:   "next middleware - success",
			method: http.MethodOptions,
			in: in{
				method: http.MethodGet,
			},
			err:      nil,
			expected: http.StatusOK,
		},
		{
			name:   "conf with method not allowed - failed",
			method: http.MethodGet,
			in: in{
				method: http.MethodGet,
			},
			err:      fmt.Errorf("method is not allowed: " + http.MethodGet),
			expected: http.StatusOK,
		},
		{
			name:   "conf with code smallest than Min - failed",
			method: http.MethodOptions,
			in: in{
				method: http.MethodOptions,
			},
			err:      fmt.Errorf("status code is smallest than minimum value: %v", cfg.StatusCodeRange.Min),
			expected: 99,
		},
		{
			name:   "conf with code biggest than Max - failed",
			method: http.MethodOptions,
			in: in{
				method: http.MethodOptions,
			},
			err:      fmt.Errorf("status code is biggest than maximum value: %v", cfg.StatusCodeRange.Max),
			expected: 600,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg.Method = tc.method
			cfg.Code = tc.expected

			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})
			handler, err := traefik_plugin_cors_preflight.New(t.Context(), next, cfg, "preflight-plugin")
			if err != nil {
				if tc.err != nil {
					assert.Equal(t, err.Error(), tc.err.Error())
					return
				}
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(t.Context(), tc.in.method, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			assert.Equal(t, recorder.Code, tc.expected)
		})
	}
}
