package main

import (
	"log"

	shipping_adapter "github.com/l-e-t-i-c-i-a/microservices/order/internal/adapters/shipping"
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

	// Shipping Adapter
	shippingAdapter, err := shipping_adapter.NewAdapter(config.GetShippingServiceUrl())
	if err != nil {
		log.Fatalf("Failed to initialize shipping stub. Error: %v", err)
	}

	// Application
	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	log.Println("Order Service is running...")
	grpcAdapter.Run()
}
