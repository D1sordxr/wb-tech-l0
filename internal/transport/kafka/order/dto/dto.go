package dto

import "time"

type Order struct {
	ID                string    `json:"order_uid" validate:"required,max=40"`
	TrackNumber       string    `json:"track_number" validate:"required,max=40"`
	Entry             string    `json:"entry" validate:"required,max=40"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,dive,required"`
	Locale            string    `json:"locale" validate:"required,bcp47_language_tag"`
	InternalSignature string    `json:"internal_signature" validate:"max=40"`
	CustomerID        string    `json:"customer_id" validate:"required,max=40"`
	DeliveryService   string    `json:"delivery_service" validate:"required,max=40"`
	ShardKey          string    `json:"shardkey" validate:"required,max=40"`
	SmID              int32     `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,max=40"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required,max=60"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required,min=5,max=8,numeric"`
	City    string `json:"city" validate:"required,min=3,max=50"`
	Address string `json:"address" validate:"required,min=3,max=100"`
	Region  string `json:"region" validate:"required,min=3,max=50"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required,max=40"`
	RequestID    string `json:"request_id" validate:"max=40"`
	Currency     string `json:"currency" validate:"required,iso4217"`
	Provider     string `json:"provider" validate:"required,max=20"`
	Amount       int32  `json:"amount" validate:"required"`
	PaymentDt    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required,max=20"`
	DeliveryCost int32  `json:"delivery_cost" validate:"required"`
	GoodsTotal   int32  `json:"goods_total" validate:"required"`
	CustomFee    int32  `json:"custom_fee"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int32  `json:"price" validate:"required"`
	RID         string `json:"rid" validate:"required,max=40"`
	Name        string `json:"name" validate:"required,max=40"`
	Sale        int32  `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int32  `json:"total_price" validate:"required"`
	NmID        int64  `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required,max=100"`
	Status      int32  `json:"status" validate:"required"`
}
