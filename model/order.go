package model

type Order struct {
	OrderUid          string   `db:"order_uid" json:"order_uid" validate:"required"`
	TrackNumber       string   `db:"track_number" json:"track_number"`
	Entry             string   `db:"entry" json:"entry"`
	Delivery          Delivery `db:"delivery" json:"delivery"`
	Payment           Payment  `db:"payment" json:"payment"`
	Items             []Items  `db:"items" json:"items"`
	Locale            string   `db:"locale" json:"locale"`
	InternalSignature string   `db:"internal_signature" json:"internal_signature"`
	CustomerId        string   `db:"customer_id" json:"customer_id"`
	DeliveryService   string   `db:"delivery_service" json:"delivery_service"`
	Shardkey          string   `db:"shardkey" json:"shardkey"`
	SmId              int      `db:"sm_id" json:"sm_id"`
	DateCreated       string   `db:"date_created" json:"date_created"`
	OofShard          string   `db:"oof_shard" json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone" validate:"e164"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Items struct {
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
