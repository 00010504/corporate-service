package models

import "database/sql"

type Shop struct {
	Id                string `json:"id"`
	Title             string `json:"title"`
	PhoneNumber       string `json:"phone_number"`
	Size              int    `json:"size"`
	Address           string `json:"address"`
	Description       string `json:"description"`
	NumberOfCashboxes int32  `json:"number_of_cashboxes"`
	CreatedAt         string `json:"created_at"`
	CreatedBy         string `json:"created_by"`
	DeletedAt         string `json:"deleted_at"`
	DeletedBy         string `json:"deleted_by"`
}

type ShortShop struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NullShortShop struct {
	ID   sql.NullString `json:"id"`
	Name sql.NullString `json:"name"`
}
