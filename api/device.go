package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PaoloModica/signing-service-challenge-go/domain"
)

type SignatureDeviceParams struct {
	Label   string                 `json:"label"`
	KeyType domain.KeyGenAlgorithm `json:"key_type"`
}

type SignatureDeviceInfoResponse struct {
	Id      string `json:"id"`
	Label   string `json:"label"`
	Counter int    `json:"counter"`
}

type SignatureDeviceInfoListResponse struct {
	Devices []SignatureDeviceInfoResponse `json:"devices"`
}

type SignatureDevicesResponse struct {
	Data SignatureDeviceInfoListResponse `json:"data"`
}

type SignatureDeviceCreationResponse struct {
	Id string `json:"id"`
}

type SignatureDeviceResponse struct {
	Data SignatureDeviceCreationResponse `json:"data"`
}

func (s *Server) HandleSignatureDeviceCreation(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	var signatureDeviceParams SignatureDeviceParams
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&signatureDeviceParams)
	if err != nil {
		WriteErrorResponse(response, http.StatusUnprocessableEntity, []string{http.StatusText(http.StatusUnprocessableEntity)})
		return
	}

	deviceId, err := s.signatureDeviceService.Create(signatureDeviceParams.Label, signatureDeviceParams.KeyType)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
	}

	WriteAPIResponse(response, http.StatusCreated, SignatureDeviceCreationResponse{Id: deviceId})
}

func (s *Server) HandleSignatureDeviceRetrieval(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{http.StatusText(http.StatusMethodNotAllowed)})
		return
	}

	deviceId := strings.TrimPrefix(request.URL.Path, "/api/v0/devices/")
	devicesList := []SignatureDeviceInfoResponse{}

	if deviceId != "" {
		device, err := s.signatureDeviceService.FindById(deviceId)
		if err != nil {
			WriteErrorResponse(response, http.StatusNotFound, []string{err.Error()})
			return
		}
		devicesList = append(devicesList, SignatureDeviceInfoResponse{Id: device.Id, Label: device.Label, Counter: device.GetSignatureCounter()})
	} else {
		devices, err := s.signatureDeviceService.FindAll()
		if err != nil {
			WriteErrorResponse(response, http.StatusInternalServerError, []string{err.Error()})
			return
		}
		for _, device := range devices {
			devicesList = append(devicesList, SignatureDeviceInfoResponse{Id: device.Id, Label: device.Label, Counter: device.GetSignatureCounter()})
		}
	}
	WriteAPIResponse(response, http.StatusOK, SignatureDeviceInfoListResponse{Devices: devicesList})
}
