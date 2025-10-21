package models

type Delivery struct {
	Name    string `db:"name"    json:"name"    validate:"required,alphaunicode_with_space,min=2,max=100"`
	Phone   string `db:"phone"   json:"phone"   validate:"required,e164"`
	Zip     string `db:"zip"     json:"zip"     validate:"required,numeric,min=5,max=8"`
	City    string `db:"city"    json:"city"    validate:"required,alphaunicode_with_space,min=3,max=50"`
	Address string `db:"address" json:"address" validate:"required,min=10,max=100"`
	Region  string `db:"region"  json:"region"  validate:"required,alphaunicode_with_space,min=3,max=60"`
	Email   string `db:"email"   json:"email"   validate:"required,email"`
}

func (d *Delivery) Validate() error {
	return validate.Struct(d)
}
