package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestHandlerFunc_JulienschmidtHttpRouter(t *testing.T) {

	tests := []struct {
		name               string
		method             string
		expectedStatusCode int
		url                string
	}{
		{
			name:   "Ok",
			method: http.MethodGet,
			// successful code of response (204) as we set it in Heartbeat method
			expectedStatusCode: http.StatusNoContent,
			url:                "/api/heartbeat",
		},
		{
			name:   "HTTP Method not allowed for this URL",
			method: http.MethodPost,
			// expected code http.StatusMethodNotAllowed because this is URL not register with this method
			expectedStatusCode: http.StatusMethodNotAllowed,
			url:                "/api/heartbeat",
		},
		{
			name:   "Requested URL not register",
			method: http.MethodGet,
			// expected code http.StatusNotFound because this is URL not register
			expectedStatusCode: http.StatusNotFound,
			url:                "/someUrl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}

			httpRouter := httprouter.New()
			h.Register(httpRouter)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, http.NoBody)

			httpRouter.ServeHTTP(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code: %d, got: %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}

// WrapServeMux wraps the http.ServeMux to implement the HandlerFunc interface.
type WrapServeMux struct {
	*http.ServeMux
}

// HandlerFunc implements the HandlerFunc interface.
func (s *WrapServeMux) HandlerFunc(method, url string, handler http.HandlerFunc) {
	s.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, req)
	})
}

func TestHandlerFunc_WrapServeMux(t *testing.T) {

	tests := []struct {
		name               string
		method             string
		expectedStatusCode int
		url                string
	}{
		{
			name:   "Ok",
			method: http.MethodGet,
			// successful code of response (204) as we set it in Heartbeat method
			expectedStatusCode: http.StatusNoContent,
			url:                "/api/heartbeat",
		},
		{
			name:   "HTTP Method not allowed for this URL",
			method: http.MethodPost,
			// expected code http.StatusMethodNotAllowed because this is URL not register with this method
			expectedStatusCode: http.StatusMethodNotAllowed,
			url:                "/api/heartbeat",
		},
		{
			name:   "Requested URL not register",
			method: http.MethodGet,
			// expected code http.StatusNotFound because this is URL not register
			expectedStatusCode: http.StatusNotFound,
			url:                "/someUrl",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}

			wrapServeMux := &WrapServeMux{
				&http.ServeMux{},
			}
			h.Register(wrapServeMux)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, http.NoBody)

			wrapServeMux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status code: %d, got: %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}
