package dtosV1

type ShopItemDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"quantity"`
	Price       float64 `json:"price"`
}
