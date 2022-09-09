package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	grpcOrderService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/entities"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConfigureMappings() error {

	// Order -> OrderDto
	err := mapper.CreateMap[*aggregate.Order, *dtos.OrderDto]()
	if err != nil {
		return err
	}

	// OrderDto -> Order
	err = mapper.CreateCustomMap[*dtos.OrderDto, *aggregate.Order](func(orderDto *dtos.OrderDto) *aggregate.Order {
		items, err := mapper.Map[[]*value_objects.ShopItem](orderDto.ShopItems)
		if err != nil {
			return nil
		}

		//payment, err := mapper.Map[*entities.Payment](orderDto.Payment)
		//if err != nil {
		//	return nil
		//}

		order, err := aggregate.NewOrder(orderDto.Id, items, orderDto.AccountEmail, orderDto.DeliveryAddress, orderDto.DeliveredTime, orderDto.CreatedAt)
		if err != nil {
			return nil
		}

		return order
	})
	if err != nil {
		return err
	}

	// ShopItem -> ShopItemDto
	err = mapper.CreateMap[*value_objects.ShopItem, *dtos.ShopItemDto]()
	if err != nil {
		return err
	}

	// ShopItemDto -> ShopItem
	err = mapper.CreateCustomMap[*dtos.ShopItemDto, *value_objects.ShopItem](func(src *dtos.ShopItemDto) *value_objects.ShopItem {
		return value_objects.CreateNewShopItem(src.Title, src.Description, src.Quantity, src.Price)
	})
	if err != nil {
		return err
	}

	// Payment -> PaymentDto
	err = mapper.CreateMap[*entities.Payment, *dtos.PaymentDto]()
	if err != nil {
		return err
	}

	// PaymentDto -> Payment
	err = mapper.CreateMap[*dtos.PaymentDto, *entities.Payment]()
	if err != nil {
		return err
	}

	// value_objects.ShopItem -> grpcOrderService.ShopItem
	err = mapper.CreateCustomMap[*value_objects.ShopItem, *grpcOrderService.ShopItem](func(src *value_objects.ShopItem) *grpcOrderService.ShopItem {
		return &grpcOrderService.ShopItem{
			Title:       src.Title(),
			Description: src.Description(),
			Quantity:    src.Quantity(),
			Price:       src.Price(),
		}
	})
	if err != nil {
		return err
	}

	// grpcOrderService.ShopItem -> value_objects.ShopItem
	err = mapper.CreateCustomMap[*grpcOrderService.ShopItem, *value_objects.ShopItem](func(src *grpcOrderService.ShopItem) *value_objects.ShopItem {
		return value_objects.CreateNewShopItem(src.Title, src.Description, src.Quantity, src.Price)
	})
	if err != nil {
		return err
	}

	// grpcOrderService.ShopItem -> dtos.ShopItemDto
	err = mapper.CreateMap[*grpcOrderService.ShopItem, *dtos.ShopItemDto]()
	if err != nil {
		return err
	}

	// grpcOrderService.Payment -> dtos.PaymentDto
	err = mapper.CreateMap[*grpcOrderService.Payment, *dtos.PaymentDto]()
	if err != nil {
		return err
	}

	//  entities.Payment -> grpcOrderService.Payment
	err = mapper.CreateMap[*entities.Payment, *grpcOrderService.Payment]()
	if err != nil {
		return err
	}

	// aggregate.Order -> grpcOrderService.Order
	err = mapper.CreateCustomMap[*aggregate.Order, *grpcOrderService.Order](func(order *aggregate.Order) *grpcOrderService.Order {
		items, err := mapper.Map[[]*grpcOrderService.ShopItem](order.ShopItems())
		if err != nil {
			return nil
		}

		payment, err := mapper.Map[*grpcOrderService.Payment](order.Payment())
		if err != nil {
			return nil
		}

		return &grpcOrderService.Order{
			OrderId:         order.Id().String(),
			DeliveryAddress: order.DeliveryAddress(),
			DeliveredTime:   timestamppb.New(order.DeliveredTime()),
			AccountEmail:    order.AccountEmail(),
			Canceled:        order.Canceled(),
			Completed:       order.Completed(),
			Paid:            order.Paid(),
			CancelReason:    order.CancelReason(),
			Submitted:       order.Submitted(),
			TotalPrice:      order.TotalPrice(),
			CreatedAt:       timestamppb.New(order.CreatedAt()),
			UpdatedAt:       timestamppb.New(order.UpdatedAt()),
			ShopItems:       items,
			Payment:         payment,
		}
	})
	if err != nil {
		return err
	}

	return nil
}
