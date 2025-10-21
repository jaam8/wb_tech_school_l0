package models

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func GenerateFakeOrder() Order {
	delivery := Delivery{
		Name:    gofakeit.Name(),
		Phone:   "+" + gofakeit.Phone(),
		Zip:     gofakeit.Zip(),
		City:    gofakeit.City(),
		Address: gofakeit.Street(),
		Region:  gofakeit.State(),
		Email:   gofakeit.Email(),
	}

	start := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 9, 30, 23, 59, 59, 0, time.UTC)
	orderedAt := gofakeit.DateRange(start, end)

	payment := Payment{
		Transaction: strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		RequestID:   ".",
		Provider:    gofakeit.Company(),
		Currency:    gofakeit.CurrencyShort(),
		PaymentDt:   int(orderedAt.Add(time.Minute).Unix()),
		Bank: gofakeit.RandomString([]string{
			"sber", "alpha", "vtb", "tinkoff",
			"yapay", "wbpay", "ozon"}),
		DeliveryCost: gofakeit.Number(0, 500),
		GoodsTotal:   gofakeit.Number(1, 10000),
		CustomFee:    gofakeit.Number(0, 10),
	}
	payment.Amount = payment.DeliveryCost + payment.GoodsTotal

	trackNumber := gofakeit.LetterN(12)
	items := make([]Item, 0)

	for j := 0; j < gofakeit.Number(1, 3); j++ {
		item := Item{
			ChrtID:      gofakeit.Number(1, 999999),
			TrackNumber: trackNumber,
			Price:       gofakeit.Number(100, 10000),
			Rid:         strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			Name:        gofakeit.ProductName(),
			Sale:        gofakeit.Number(0, 99),
			Size:        gofakeit.Numerify("#"),
			NmID:        gofakeit.Number(1, 9999999),
			Brand:       gofakeit.Company(),
			Status:      gofakeit.Number(1, 5),
		}
		item.TotalPrice = int(float64(item.Price) * (1 - float64(item.Sale)/100.0))

		items = append(items, item)
	}

	order := Order{
		OrderUID:          strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		TrackNumber:       trackNumber,
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            strings.ToUpper(gofakeit.CountryAbr()),
		InternalSignature: ".",
		CustomerID:        strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		DeliveryService:   gofakeit.Company(),
		Shardkey:          gofakeit.Numerify("#"),
		SmID:              gofakeit.Number(1, 100),
		DateCreated:       orderedAt,
		OofShard:          gofakeit.Numerify("#"),
	}

	return order
}
