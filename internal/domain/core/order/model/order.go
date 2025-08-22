package model

import "time"

type Order struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          Delivery
	Payment           Payment
	Items             []Item
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SmID              int32
	DateCreated       time.Time
	OofShard          string
}

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Payment struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int32
	PaymentDt    int64
	Bank         string
	DeliveryCost int32
	GoodsTotal   int32
	CustomFee    int32
}

type Item struct {
	ChartID     int64
	TrackNumber string
	Price       int32
	RID         string
	Name        string
	Sale        int32
	Size        string
	TotalPrice  int32
	NmID        int64
	Brand       string
	Status      int32
}
