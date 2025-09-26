package models

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction" validate:"required,min=10,max=60"`
	RequestId    string `json:"request_id" db:"request_id" validate:"required"`
	Provider     string `json:"provider" db:"provider" validate:"required,min=2,max=25"`
	Currency     string `json:"currency" db:"currency" validate:"required,iso4217"`
	Amount       int    `json:"amount" db:"amount" validate:"required,gt=0"`
	PaymentDt    int    `json:"payment_dt" db:"payment_dt" validate:"required,datetime"`
	Bank         string `json:"bank" db:"bank" validate:"required,oneof=sber alpha vtb tinkoff"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost" validate:"required,gte=0"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total" validate:"required,gt=0"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee" validate:"required,gte=0"`
}

func (p *Payment) Validate() error {
	return validate.Struct(p)
}
