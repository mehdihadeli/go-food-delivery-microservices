package mappings

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	grpcOrderService "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/entities"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/value_objects"
)

func ConfigureMappings() error {
	if err := mapper.CreateCustomMap[*value_objects.ShopItem, *grpcOrderService.ShopItem](func(src *value_objects.ShopItem) *grpcOrderService.ShopItem {
		return &grpcOrderService.ShopItem{
			Title:       src.Title,
			Description: src.Description,
			Quantity:    src.Quantity,
			Price:       src.Price,
		}
	}); err != nil {
		return err
	}

	if err := mapper.CreateCustomMap[*grpcOrderService.ShopItem, *value_objects.ShopItem](func(src *grpcOrderService.ShopItem) *value_objects.ShopItem {
		return &value_objects.ShopItem{
			Title:       src.Title,
			Description: src.Description,
			Quantity:    src.Quantity,
			Price:       src.Price,
		}
	}); err != nil {
		return err
	}

	// ShopItem -> ShopItemDto
	err := mapper.CreateMap[*value_objects.ShopItem, *dtos.ShopItemDto]()
	//err := mapper.CreateMap[value_objects.ShopItem, dtos.ShopItemDto]()
	if err != nil {
		return err
	}

	// ShopItemDto -> ShopItem
	err = mapper.CreateMap[*dtos.ShopItemDto, *value_objects.ShopItem]()
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

	// Order -> OrderDto
	err = mapper.CreateMap[*aggregate.Order, *dtos.OrderDto]()
	if err != nil {
		return err
	}

	return nil
}
