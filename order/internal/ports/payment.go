package ports

import "github.com/niryb/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(domain.Order) error //Order ou OrderService?
}

//Veja os outros arquivos de código para definição de portas já existentes no projeto para se basear
//em como fazer a implementação do arquivo internal/ports/payment.go.
