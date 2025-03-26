package dto

type SearchProduct struct {
	Description string  `json:"description" binding:"required"`
	PriceMin    float64 `json:"price_min"`
	PriceMax    float64 `json:"price_max"`
}
