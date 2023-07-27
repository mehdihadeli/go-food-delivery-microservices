package dtosV1

type ShopItemReadDto struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Quantity    uint64  `json:"quantity"`
	Price       float64 `json:"price"`
}
