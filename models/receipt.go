package models

type Receipt struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Text         string `json:"text"`
	Logo         bool   `json:"logo"`
	Date         bool   `json:"date"`
	Shop         bool   `json:"shop"`
	WorkingHours bool   `json:"working_hours"`
	Seller       bool   `json:"seller"`
	Cashier      bool   `json:"cashier"`
	Customer     bool   `json:"customer"`
	Contact      bool   `json:"contact"`
	INN          bool   `json:"inn"`
	Barcode      bool   `json:"barcode"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	DeletedAt    string `json:"deleted_at"`
	DeletedBy    string `json:"deleted_by"`
}
