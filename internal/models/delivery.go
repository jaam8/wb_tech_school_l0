package models

type Delivery struct {
	Name    string `json:"name" db:"name" validate:"required,alphaunicode,min=2,max=100"`
	Phone   string `json:"phone" db:"phone" validate:"required,e164"`
	Zip     string `json:"zip" db:"zip" validate:"required,postcode_iso3166_alpha2_field=Country"`
	City    string `json:"city" db:"city" validate:"required,alphaunicode,min=3,max=50"`
	Address string `json:"address" db:"address" validate:"required,min=10,max=100"`
	Region  string `json:"region" db:"region" validate:"required,alphaunicode,min=3,max=60"`
	Email   string `json:"email" db:"email" validate:"required,email"`
}

func (d *Delivery) Validate() error {
	return validate.Struct(d)
}
