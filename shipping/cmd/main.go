package main

import (
	"log"

	"github.com/niryb/microservices/shipping/config"
	"github.com/niryb/microservices/shipping/internal/adapters/grpc"
	"github.com/niryb/microservices/shipping/internal/application/core/api"
)

func main() {
	application := api.NewApplication()

	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())

	log.Printf("Starting Shipping Service on port %d...", config.GetApplicationPort())
	grpcAdapter.Run()
}
