package api

import (
	"github.com/niryb/microservices/shipping/internal/application/core/domain"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a Application) CreateShipping(shipping domain.Shipping) (domain.Shipping, error) {
	shipping.CalculateDelivery()

	return shipping, nil
}
