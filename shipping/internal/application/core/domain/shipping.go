package domain

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

type Shipping struct {
	OrderID      int64
	Items        []ShippingItem
	DeliveryDays int32
}

func NewShipping(orderId int64, items []ShippingItem) Shipping {
	shipping := Shipping{
		OrderID: orderId,
		Items:   items,
	}

	shipping.CalculateDelivery()

	return shipping
}

func (s *Shipping) CalculateDelivery() {
	totalQty := int32(0)
	for _, item := range s.Items {
		totalQty += item.Quantity
	}
	s.DeliveryDays = 1 + (totalQty / 5)
}
