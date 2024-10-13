package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/api"
)

func TestServer(t *testing.T) {
	baseUrl := "http://localhost:8080"
	testSetup(t, baseUrl)

	t.Run("GET /api/v0/health returns 200 OK", func(t *testing.T) {
		http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v0/health", baseUrl), nil)
		response := httptest.NewRecorder()

		assertResponseStatusCode(t, http.StatusOK, response.Result().StatusCode)
	})
}

func testSetup(t *testing.T, baseUrl string) {
	t.Helper()

	server := api.NewServer(baseUrl)
	server.Run()
}

func assertResponseStatusCode(t *testing.T, expected int, got int) {
	t.Helper()

	if expected != got {
		t.Errorf("expected %d HTTP status code, got %d", expected, got)
	}
}
