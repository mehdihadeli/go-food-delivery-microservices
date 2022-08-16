package value_objects

import (
	"fmt"
)

type ShopItem struct {
	title       string
	description string
	quantity    uint64
	price       float64
}

func CreateNewShopItem(title string, description string, quantity uint64, price float64) *ShopItem {
	return &ShopItem{
		title:       title,
		description: description,
		quantity:    quantity,
		price:       price,
	}
}

func (s *ShopItem) Title() string {
	return s.title
}

func (s *ShopItem) Description() string {
	return s.description
}

func (s *ShopItem) Quantity() uint64 {
	return s.quantity
}

func (s *ShopItem) Price() float64 {
	return s.price
}

func (s *ShopItem) String() string {
	return fmt.Sprintf("Title: {%s}, Description: {%s}, Quantity: {%v}, Price: {%v},",
		s.title,
		s.description,
		s.quantity,
		s.price,
	)
}
