package value_objects

import (
	"fmt"
)

type ShopItem struct {
	Title       string  `json:"title" bson:"title,omitempty"`
	Description string  `json:"description" bson:"description,omitempty"`
	Quantity    uint64  `json:"quantity" bson:"quantity,omitempty"`
	Price       float64 `json:"price" bson:"price,omitempty"`
}

func (s *ShopItem) String() string {
	return fmt.Sprintf("Title: {%s}, Description: {%s}, Quantity: {%v}, Price: {%v},",
		s.Title,
		s.Description,
		s.Quantity,
		s.Price,
	)
}
