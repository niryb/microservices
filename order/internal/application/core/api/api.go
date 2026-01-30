package api

import (
	"log"

	"github.com/niryb/microservices/order/internal/application/core/domain"
	"github.com/niryb/microservices/order/internal/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Verifica se a quantidade total de itens excede 50
	if order.TotalQuantity() > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, "Total quantity of items cannot exceed 50")
	}

	if err := a.db.CheckProductsExist(order.OrderItems); err != nil {
		return domain.Order{}, status.Errorf(codes.NotFound, err.Error())
	}

	// Salva o pedido com status "Pending"
	err := a.db.Save(&order)

	// Se erro ao salvar, retorna
	if err != nil {
		return domain.Order{}, err
	}

	// Tenta cobrar o pagamento
	paymentErr := a.payment.Charge(order)
	if paymentErr != nil {
		// Se falhar, atualiza o status do pedido para "Canceled" e retorna o erro
		order.Status = "Canceled"
		a.db.Update(order)
		return domain.Order{}, paymentErr // Retorna o erro de pagamento
	}

	shippingErr := a.shipping.ShipOrder(order)
	if shippingErr != nil {
		log.Printf("Erro ao solicitar envio: %v", shippingErr)
		order.Status = "Paid"
	} else {
		order.Status = "Shipped"
	}
	a.db.Update(order)

	// Retorna o pedido atualizado
	return order, nil
}
