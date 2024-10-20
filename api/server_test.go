package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaoloModica/signing-service-challenge-go/api"
	"github.com/PaoloModica/signing-service-challenge-go/domain"
	test_utils "github.com/PaoloModica/signing-service-challenge-go/internal"
)

func TestServer(t *testing.T) {
	baseUrl := "http://localhost:8080"

	device, _ := domain.NewSignatureDevice("testDevice1", []byte("privateKey"), "RSA")
	store := test_utils.StubSignatureDeviceStore{
		Store: map[string]*domain.SignatureDevice{device.Id: device},
	}
	repository, _ := domain.NewSignatureDeviceRepository(&store)
	service, _ := domain.NewSignatureDeviceService(repository)

	server := api.NewServer(baseUrl, service)
	server.InitializeRouter()

	t.Run("GET /api/v0/unknown returns 404 Not Found", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v0/unknown", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, http.StatusNotFound, response.Result().StatusCode)
	})
	t.Run("GET /api/v0/health returns 200 OK", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v0/health", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, http.StatusOK, response.Result().StatusCode)
	})
	t.Run("POST /api/v0/devices returns 201 Created", func(t *testing.T) {
		signatureDeviceParam := api.SignatureDeviceParams{
			Label:   "testDevice",
			KeyType: "RSA",
		}
		marshalledDeviceParam, _ := json.Marshal(signatureDeviceParam)
		request, _ := http.NewRequest(http.MethodPost, "/api/v0/devices", bytes.NewReader(marshalledDeviceParam))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseResult := response.Result()
		assertResponseStatusCode(t, http.StatusCreated, responseResult.StatusCode)

		defer responseResult.Body.Close()

		var deviceCreationResponse api.SignatureDeviceResponse
		json.NewDecoder(responseResult.Body).Decode(&deviceCreationResponse)

		if deviceCreationResponse.Data.Id == "" {
			t.Errorf("expected signature device creation response to return device ID")
		}
	})
	t.Run("PUT /api/v0/devices returns 405 Method Not Allowed", func(t *testing.T) {
		signatureDeviceParam := api.SignatureDeviceParams{
			Label:   "testDevice",
			KeyType: "RSA",
		}
		marshalledDeviceParam, _ := json.Marshal(signatureDeviceParam)
		request, _ := http.NewRequest(http.MethodPut, "/api/v0/devices", bytes.NewReader(marshalledDeviceParam))
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseResult := response.Result()
		assertResponseStatusCode(t, http.StatusMethodNotAllowed, responseResult.StatusCode)
	})
	t.Run("GET /api/v0/devices/ returns 200 and list of devices", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v0/devices/", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseResult := response.Result()
		assertResponseStatusCode(t, http.StatusOK, responseResult.StatusCode)

		expectedDeviceLen := 2
		var devicesResponse api.SignatureDevicesResponse

		defer responseResult.Body.Close()
		responseResultBody, _ := io.ReadAll(responseResult.Body)
		json.Unmarshal(responseResultBody, &devicesResponse)

		test_utils.AssertSignatureDeviceStoreLen(t, expectedDeviceLen, len(devicesResponse.Data.Devices))
	})
	t.Run("GET /api/v0/devices/:id returns 200 and list of devices with single device", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v0/devices/%s", device.Id), nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseResult := response.Result()
		assertResponseStatusCode(t, http.StatusOK, responseResult.StatusCode)

		expectedDeviceLen := 1
		var devicesResponse api.SignatureDevicesResponse

		defer responseResult.Body.Close()
		responseResultBody, _ := io.ReadAll(responseResult.Body)
		json.Unmarshal(responseResultBody, &devicesResponse)

		test_utils.AssertSignatureDeviceStoreLen(t, expectedDeviceLen, len(devicesResponse.Data.Devices))

		gotDevice := devicesResponse.Data.Devices[0]
		if gotDevice.Id != device.Id {
			t.Errorf("expected signature device retrieve to be %s, got %s", device.Id, gotDevice.Id)
		}
	})
	t.Run("GET /api/v0/devices/:id returns 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/v0/devices/unknown", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		responseResult := response.Result()
		assertResponseStatusCode(t, http.StatusNotFound, responseResult.StatusCode)
	})
}

func assertResponseStatusCode(t *testing.T, expected int, got int) {
	t.Helper()

	if expected != got {
		t.Errorf("expected %d HTTP status code, got %d", expected, got)
	}
}
