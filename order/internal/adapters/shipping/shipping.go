package shipping_adapter

import (
	"context"
	"fmt"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	"github.com/niryb/microservices-proto/golang/shipping"

	"github.com/niryb/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1 * time.Second)),
		grpc_retry.WithMax(5),
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
	}

	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(shippingServiceUrl, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping service: %v", err)
	}

	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

// ShipOrder envia o pedido para o serviço de Shipping
func (a *Adapter) ShipOrder(order domain.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var protoItems []*shipping.ShippingItem
	for _, item := range order.OrderItems {
		protoItems = append(protoItems, &shipping.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	// Chama o serviço de Shipping
	_, err := a.shipping.Create(ctx, &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items:   protoItems, // Passamos a lista convertida
	})

	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("TIMEOUT EXCEDIDO: O serviço de Shipping demorou mais de 2s para o pedido %d", order.ID)
		}
		return err
	}

	return nil
}
