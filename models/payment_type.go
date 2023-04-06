package models

type PaymentType struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at"`
}
