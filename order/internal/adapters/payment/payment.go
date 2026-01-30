package payment

import (
	"context"

	//"log"
	"fmt"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/niryb/microservices-proto/golang/payment"
	"github.com/niryb/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption

	// Configura as opções de retry
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1 * time.Second)),
		grpc_retry.WithMax(5),
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
	}

	// Adiciona o interceptor de retry às opções do gRPC
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order domain.Order) error {

	// Configura um contexto com timeout para a chamada gRPC
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	// Chama o método Create do serviço de pagamento
	_, err := a.payment.Create(ctx, &payment.CreatePaymentRequest{

		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})
	// Verifica se houve erro na chamada gRPC
	if err != nil {
		// Verifica se o erro foi por timeout
		if status.Code(err) == codes.DeadlineExceeded {
			// Log da pratica
			log.Printf("---- TIMEOUT EXCEDIDO ---- : Chamada ao serviço de Payment demorou mais de 2 segundos para o pedido %d", order.ID)
		}
		return err
	}

	return nil
}
