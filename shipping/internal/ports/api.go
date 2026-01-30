package ports

import "github.com/niryb/microservices/shipping/internal/application/core/domain"

type APIPort interface {
	CreateShipping(shipping domain.Shipping) (domain.Shipping, error)
}
