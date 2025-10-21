package models

type Payment struct {
	Transaction  string `db:"transaction"   json:"transaction"   validate:"required,min=10,max=60"`
	RequestID    string `db:"request_id"    json:"request_id"    validate:"required"`
	Provider     string `db:"provider"      json:"provider"      validate:"required,min=2,max=25"`
	Currency     string `db:"currency"      json:"currency"      validate:"required,iso4217"`
	Amount       int    `db:"amount"        json:"amount"        validate:"required,gt=0"`
	PaymentDt    int    `db:"payment_dt"    json:"payment_dt"    validate:"required"`
	Bank         string `db:"bank"          json:"bank"          validate:"required,oneof=sber alpha vtb tinkoff"`
	DeliveryCost int    `db:"delivery_cost" json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `db:"goods_total"   json:"goods_total"   validate:"required,gt=0"`
	CustomFee    int    `db:"custom_fee"    json:"custom_fee"    validate:"gte=0,lte=99"`
}

func (p *Payment) Validate() error {
	return validate.Struct(p)
}
