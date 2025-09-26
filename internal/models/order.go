package models

import "time"

type Order struct {
	OrderUid          string    `json:"order_uid" db:"order_uid" validate:"required,min=10,max=60"`
	TrackNumber       string    `json:"track_number" db:"track_number" validate:"required,min=10,max=60"`
	Entry             string    `json:"entry" db:"entry" validate:"required,min=3,max=10"`
	Delivery          Delivery  `json:"delivery" db:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" db:"payment" validate:"required"`
	Items             []Item    `json:"items" db:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" db:"locale" validate:"required,iso3166_1_alpha2"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature" validate:"required"`
	CustomerId        string    `json:"customer_id" db:"customer_id" validate:"required,min=2,max=50"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" db:"shardkey" validate:"required,min=1,max=10"`
	SmId              int       `json:"sm_id" db:"sm_id" validate:"required,gt=0"`
	DateCreated       time.Time `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" db:"oof_shard" validate:"required,min=1,max=10"`
}

func (o *Order) Validate() error {
	return validate.Struct(o)
}
