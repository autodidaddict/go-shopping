package redis

type redisProduct struct {
	SKU          string `redis:"sku"`
	Name         string `redis:"name"`
	Description  string `redis:"description"`
	Manufacturer string `redis:"mfr"`
	Model        string `redis:"model"`
	Price        int64  `redis:"price"`
}

type redisCategory struct {
	Name        string `redis:"name"`
	Description string `redis:"description"`
}
