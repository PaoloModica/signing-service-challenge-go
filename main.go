package main

import (
	"log"

	"github.com/PaoloModica/signing-service-challenge-go/api"
	"github.com/PaoloModica/signing-service-challenge-go/domain"
	"github.com/PaoloModica/signing-service-challenge-go/persistence"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	signatureDeviceInMemoryStore, err := persistence.NewInMemorySignatureDeviceStore()
	if err != nil {
		log.Fatalf("an error occurred while setting signature device store: %s", err.Error())
		return
	}
	signatureDeviceRepository, err := domain.NewSignatureDeviceRepository(signatureDeviceInMemoryStore)
	if err != nil {
		log.Fatalf("an error occurred while setting signature device repository: %s", err.Error())
		return
	}
	signatureDeviceService, err := domain.NewSignatureDeviceService(signatureDeviceRepository)
	if err != nil {
		log.Fatalf("an error occurred while setting signature device service: %s", err.Error())
		return
	}

	server := api.NewServer(ListenAddress, signatureDeviceService)
	server.InitializeRouter()

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
