package prober

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpProber_Check(t *testing.T) {
	tests := []struct {
		name           string
		handlerStatus  int
		expectedResult bool
	}{
		{
			name:           "Success 200",
			handlerStatus:  http.StatusOK,
			expectedResult: true,
		},
		{
			name:           "Failure 500",
			handlerStatus:  http.StatusInternalServerError,
			expectedResult: false,
		},
		{
			name:           "Failure 404",
			handlerStatus:  http.StatusNotFound,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
			}))
			defer server.Close()

			prober := NewHttpProber(server.URL)
			result := prober.Check()

			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
