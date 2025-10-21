package models

type Item struct {
	ChrtID      int    `db:"chrt_id"      json:"chrt_id"      validate:"required,gt=0"`
	TrackNumber string `db:"track_number" json:"track_number" validate:"required,min=10,max=50"`
	Price       int    `db:"price"        json:"price"        validate:"required,gt=0"`
	Rid         string `db:"rid"          json:"rid"          validate:"required,hexadecimal,min=10,max=50"`
	Name        string `db:"name"         json:"name"         validate:"required,min=3,max=100"`
	Sale        int    `db:"sale"         json:"sale"         validate:"gte=0,lte=99"`
	Size        string `db:"size"         json:"size"         validate:"required,min=1"`
	TotalPrice  int    `db:"total_price"  json:"total_price"  validate:"required,gt=0"`
	NmID        int    `db:"nm_id"        json:"nm_id"        validate:"required,gt=0"`
	Brand       string `db:"brand"        json:"brand"        validate:"required,min=3,max=50"`
	Status      int    `db:"status"       json:"status"       validate:"required,gt=0"`
}

func (i *Item) Validate() error {
	return validate.Struct(i)
}
