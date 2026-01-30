package ports

import "github.com/niryb/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	Update(order domain.Order) error
	CheckProductsExist(orderItems []domain.OrderItem) error
}
