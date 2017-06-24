package service

type productDetails struct {
	SKU            string `json:"sku"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	Price          int64  `json:"price"`
	StockRemaining uint32 `json:"stock_remaining"`
}
