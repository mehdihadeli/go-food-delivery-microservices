package mappings

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	grpcOrderService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts/proto/service_clients"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/models/orders/read_models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/models/orders/value_objects"
)

func ConfigureOrdersMappings() error {
	// Order -> OrderDto
	err := mapper.CreateMap[*aggregate.Order, *dtosV1.OrderDto]()
	if err != nil {
		return err
	}

	// OrderDto -> Order
	err = mapper.CreateCustomMap[*dtosV1.OrderDto, *aggregate.Order](func(orderDto *dtosV1.OrderDto) *aggregate.Order {
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

	// read_models.OrderReadModel -> dtos.OrderReadDto
	err = mapper.CreateMap[*read_models.OrderReadModel, *dtosV1.OrderReadDto]()
	if err != nil {
		return err
	}

	// dtos.OrderReadDto -> grpcOrderService.OrderReadModel
	// custom filed map not support yet like ForMember so we have to create a custom map because of some timestamp fields map to time.Time
	err = mapper.CreateCustomMap[*dtosV1.OrderReadDto, *grpcOrderService.OrderReadModel](func(orderReadDto *dtosV1.OrderReadDto) *grpcOrderService.OrderReadModel {
		if orderReadDto == nil {
			return nil
		}
		items, err := mapper.Map[[]*grpcOrderService.ShopItemReadModel](orderReadDto.ShopItems)
		if err != nil {
			return nil
		}

		return &grpcOrderService.OrderReadModel{
			Id:              orderReadDto.Id,
			OrderId:         orderReadDto.OrderId,
			PaymentId:       orderReadDto.PaymentId,
			DeliveredTime:   timestamppb.New(orderReadDto.DeliveredTime),
			TotalPrice:      orderReadDto.TotalPrice,
			DeliveryAddress: orderReadDto.DeliveryAddress,
			AccountEmail:    orderReadDto.AccountEmail,
			Canceled:        orderReadDto.Canceled,
			Completed:       orderReadDto.Completed,
			Paid:            orderReadDto.Paid,
			Submitted:       orderReadDto.Submitted,
			CancelReason:    orderReadDto.CancelReason,
			ShopItems:       items,
			CreatedAt:       timestamppb.New(orderReadDto.CreatedAt),
			UpdatedAt:       timestamppb.New(orderReadDto.UpdatedAt),
		}
	})
	if err != nil {
		return err
	}

	// dtos.ShopItemReadDto -> grpcOrderService.ShopItemReadModel
	err = mapper.CreateMap[*dtosV1.ShopItemReadDto, *grpcOrderService.ShopItemReadModel]()
	if err != nil {
		return err
	}

	// ShopItem -> ShopItemDto
	err = mapper.CreateMap[*value_objects.ShopItem, *dtosV1.ShopItemDto]()
	if err != nil {
		return err
	}

	// ShopItemDto -> ShopItem
	err = mapper.CreateCustomMap[*dtosV1.ShopItemDto, *value_objects.ShopItem](func(src *dtosV1.ShopItemDto) *value_objects.ShopItem {
		return value_objects.CreateNewShopItem(src.Title, src.Description, src.Quantity, src.Price)
	})
	if err != nil {
		return err
	}

	// dtos.ShopItemDto -> read_models.ShopItemReadModel
	err = mapper.CreateMap[*dtosV1.ShopItemDto, *read_models.ShopItemReadModel]()
	if err != nil {
		return err
	}

	// read_models.ShopItemReadModel -> dtos.ShopItemReadDto
	err = mapper.CreateMap[*read_models.ShopItemReadModel, *dtosV1.ShopItemReadDto]()
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
	err = mapper.CreateMap[*grpcOrderService.ShopItem, *dtosV1.ShopItemDto]()
	if err != nil {
		return err
	}

	// aggregate.Order -> grpcOrderService.Order
	err = mapper.CreateCustomMap[*aggregate.Order, *grpcOrderService.Order](func(order *aggregate.Order) *grpcOrderService.Order {
		items, err := mapper.Map[[]*grpcOrderService.ShopItem](order.ShopItems())
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
			PaymentId:       order.PaymentId().String(),
		}
	})
	if err != nil {
		return err
	}

	err = mapper.CreateCustomMap[*utils.ListResult[*dtosV1.OrderReadDto], *grpcOrderService.GetOrdersRes](func(orders *utils.ListResult[*dtosV1.OrderReadDto]) *grpcOrderService.GetOrdersRes {
		o, err := mapper.Map[[]*grpcOrderService.OrderReadModel](orders.Items)
		if err != nil {
			return nil
		}
		return &grpcOrderService.GetOrdersRes{
			Pagination: &grpcOrderService.Pagination{
				Size:       int32(orders.Size),
				Page:       int32(orders.Page),
				TotalItems: orders.TotalItems,
				TotalPages: int32(orders.TotalPage),
			},
			Orders: o,
		}
	})
	if err != nil {
		return err
	}

	return nil
}
