package mock

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
	"wb-tech-l0/internal/transport/kafka/order/dto"
)

type Generator struct{}

func NewMockGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateOrder() dto.Order {
	orderID := g.generateOrderID()
	trackNumber := g.generateTrackNumber()

	return dto.Order{
		ID:                orderID,
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Delivery:          g.GenerateDelivery(),
		Payment:           g.GeneratePayment(orderID),
		Items:             g.GenerateItems(trackNumber),
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        g.generateRandomString(10),
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}

func (g *Generator) GenerateDelivery() dto.Delivery {
	return dto.Delivery{
		Name:    g.generateName(),
		Phone:   g.generatePhone(),
		Zip:     g.generateZip(),
		City:    g.generateCity(),
		Address: g.generateAddress(),
		Region:  g.generateRegion(),
		Email:   g.generateEmail(),
	}
}

func (g *Generator) GeneratePayment(transactionID string) dto.Payment {
	return dto.Payment{
		Transaction:  transactionID,
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       g.generateAmount(1000, 5000),
		PaymentDt:    time.Now().Unix(),
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   g.generateAmount(300, 1000),
		CustomFee:    0,
	}
}

func (g *Generator) GenerateItems(trackNumber string) []dto.Item {
	count := g.randomInt(1, 5)
	items := make([]dto.Item, count)

	for i := range items {
		items[i] = dto.Item{
			ChrtID:      g.randomInt64(9000000, 9999999),
			TrackNumber: trackNumber,
			Price:       g.generateAmount(100, 1000),
			RID:         g.generateRandomString(20),
			Name:        g.generateProductName(),
			Sale:        g.randomInt32(0, 50),
			Size:        "0",
			TotalPrice:  g.generateAmount(50, 500),
			NmID:        g.randomInt64(2000000, 2999999),
			Brand:       g.generateBrand(),
			Status:      202,
		}
	}

	return items
}

func (g *Generator) GenerateMultipleOrders(count int) []dto.Order {
	orders := make([]dto.Order, count)
	for i := range orders {
		orders[i] = g.GenerateOrder()
		time.Sleep(1 * time.Millisecond)
	}
	return orders
}

func (g *Generator) generateOrderID() string {
	return g.generateRandomString(16) + "test"
}

func (g *Generator) generateTrackNumber() string {
	return "WBILM" + g.generateRandomString(8) + "TRACK"
}

func (g *Generator) generateName() string {
	firstNames := []string{"John", "Jane", "Alex", "Maria", "David", "Sarah", "Mike", "Anna"}
	lastNames := []string{"Smith", "Johnson", "Brown", "Davis", "Wilson", "Taylor", "Clark", "Walker"}

	return g.randomChoice(firstNames) + " " + g.randomChoice(lastNames)
}

func (g *Generator) generatePhone() string {
	return "+7" + g.generateNumericString(10)
}

func (g *Generator) generateZip() string {
	return g.generateNumericString(6)
}

func (g *Generator) generateCity() string {
	cities := []string{"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg", "Kazan", "Nizhny Novgorod", "Chelyabinsk", "Samara"}
	return g.randomChoice(cities)
}

func (g *Generator) generateAddress() string {
	streets := []string{"Lenina", "Gorkogo", "Pushkina", "Lermontova", "Sovetskaya", "Centralnaya", "Molodezhnaya", "Shkolnaya"}
	return fmt.Sprintf("%s st., %d", g.randomChoice(streets), g.randomInt(1, 100))
}

func (g *Generator) generateRegion() string {
	regions := []string{"Moscow Oblast", "Leningrad Oblast", "Sverdlovsk Oblast", "Republic of Tatarstan", "Krasnodar Krai"}
	return g.randomChoice(regions)
}

func (g *Generator) generateEmail() string {
	domains := []string{"gmail.com", "yahoo.com", "mail.ru", "yandex.ru"}
	return fmt.Sprintf("%s@%s", g.generateRandomString(8), g.randomChoice(domains))
}

func (g *Generator) generateProductName() string {
	products := []string{
		"Smartphone", "Laptop", "Headphones", "Keyboard", "Mouse",
		"Monitor", "Tablet", "Smartwatch", "Camera", "Printer",
		"Mascaras", "Lipstick", "Foundation", "Eyeshadow", "Perfume",
	}
	return g.randomChoice(products)
}

func (g *Generator) generateBrand() string {
	brands := []string{
		"Samsung", "Apple", "Sony", "LG", "Xiaomi",
		"Huawei", "Lenovo", "Dell", "HP", "Canon",
		"Vivienne Sabo", "L'Oreal", "Maybelline", "MAC", "Chanel",
	}
	return g.randomChoice(brands)
}

func (g *Generator) generateAmount(min, max int32) int32 {
	return min + g.randomInt32(0, max-min)
}

func (g *Generator) generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}

	return string(result)
}

func (g *Generator) generateNumericString(length int) string {
	const digits = "0123456789"
	result := make([]byte, length)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		result[i] = digits[num.Int64()]
	}

	return string(result)
}

func (g *Generator) randomChoice(options []string) string {
	if len(options) == 0 {
		return ""
	}
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	return options[num.Int64()]
}

func (g *Generator) randomInt(min, max int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int(num.Int64())
}

func (g *Generator) randomInt32(min, max int32) int32 {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int32(num.Int64())
}

func (g *Generator) randomInt64(min, max int64) int64 {
	num, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + num.Int64()
}
