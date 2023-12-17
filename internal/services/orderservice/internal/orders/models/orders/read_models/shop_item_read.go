package read_models

type ShopItemReadModel struct {
	Title       string  `json:"title,omitempty"       bson:"title,omitempty"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Quantity    uint64  `json:"quantity,omitempty"    bson:"quantity,omitempty"`
	Price       float64 `json:"price,omitempty"       bson:"price,omitempty"`
}

func NewShopItemReadModel(title string, description string, quantity uint64, price float64) *ShopItemReadModel {
	return &ShopItemReadModel{Title: title, Description: description, Quantity: quantity, Price: price}
}
