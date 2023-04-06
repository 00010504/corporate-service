package models

type Cashbox struct {
	Id             string   `json:"id"`
	ShopId         string   `json:"shop_id"`
	Title          string   `json:"title"`
	PaymentTypeIds []string `json:"payment_type_ids"`
	ReceiptTypeId  string   `json:"receipt_type_id"`
	CreatedAt      string   `json:"created_at"`
	CreatedBy      string   `json:"created_by"`
	DeletedAt      string   `json:"deleted_at"`
	DeletedBy      string   `json:"deleted_by"`
}
