package models

type Item struct {
	ChrtId      int    `json:"chrt_id" db:"chrt_id" validate:"required,gt=0"`
	TrackNumber string `json:"track_number" db:"track_number" validate:"required,min=10,max=50"`
	Price       int    `json:"price" db:"price" validate:"required,gt=0"`
	Rid         string `json:"rid" db:"rid" validate:"required,hexadecimal,min=10,max=50"`
	Name        string `json:"name" db:"name" validate:"required,min=3,max=100"`
	Sale        int    `json:"sale" db:"sale" validate:"gte=0,lte=99"`
	Size        string `json:"size" db:"size" validate:"required,min=1"`
	TotalPrice  int    `json:"total_price" db:"total_price" validate:"required,gt=0"`
	NmId        int    `json:"nm_id" db:"nm_id" validate:"required,gt=0"`
	Brand       string `json:"brand" db:"brand" validate:"required,min=3,max=50"`
	Status      int    `json:"status" db:"status" validate:"required,gt=0"`
}

func (i *Item) Validate() error {
	return validate.Struct(i)
}
