package main

import (
	"log"

	"github.com/niryb/microservices/order/config"
	"github.com/niryb/microservices/order/internal/adapters/db"
	"github.com/niryb/microservices/order/internal/adapters/grpc"
	"github.com/niryb/microservices/order/internal/adapters/payment"
	"github.com/niryb/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	log.Println("Order service starting...")
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}
	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
