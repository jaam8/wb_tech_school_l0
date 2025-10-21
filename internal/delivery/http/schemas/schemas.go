package schemas

type ErrorResponse struct {
	Success bool   `example:"false"         json:"success"`
	Error   string `example:"error message" json:"error"`
}
