package ports

import "github.com/niryb/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	ShipOrder(order domain.Order) error
}
