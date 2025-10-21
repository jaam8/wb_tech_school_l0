package models

import (
	"time"
)

type Order struct {
	OrderUID          string    `db:"order_uid"          json:"order_uid"          validate:"required,hexadecimal,min=10,max=60"`
	TrackNumber       string    `db:"track_number"       json:"track_number"       validate:"required,min=10,max=60"`
	Entry             string    `db:"entry"              json:"entry"              validate:"required,min=3,max=10"`
	Delivery          Delivery  `db:"delivery"           json:"delivery"           validate:"required"`
	Payment           Payment   `db:"payment"            json:"payment"            validate:"required"`
	Items             []Item    `db:"items"              json:"items"              validate:"required,min=1,dive"`
	Locale            string    `db:"locale"             json:"locale"             validate:"required,country_code"`
	InternalSignature string    `db:"internal_signature" json:"internal_signature" validate:"required"`
	CustomerID        string    `db:"customer_id"        json:"customer_id"        validate:"required,min=2,max=50"`
	DeliveryService   string    `db:"delivery_service"   json:"delivery_service"   validate:"required"`
	Shardkey          string    `db:"shardkey"           json:"shardkey"           validate:"required,min=1,max=10"`
	SmID              int       `db:"sm_id"              json:"sm_id"              validate:"required,gt=0"`
	DateCreated       time.Time `db:"date_created"       json:"date_created"       validate:"required"`
	OofShard          string    `db:"oof_shard"          json:"oof_shard"          validate:"required,min=1,max=10"`
}

func (o *Order) Validate() error {
	return validate.Struct(o)
}
