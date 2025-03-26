package dto

type SearchProduct struct {
	Description string  `json:"description" binding:"omitempty"`
	PriceMin    float64 `json:"price_min" binding:"omitempty"`
	PriceMax    float64 `json:"price_max" binding:"omitempty"`
}
