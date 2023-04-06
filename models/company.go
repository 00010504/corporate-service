package models

import (
	"database/sql"
	"genproto/common"
	"genproto/corporate_service"
)

type Company struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	CompanyTypeId string `json:"company_type_id"`
	CreatedBy     string `json:"created_by"`
	DeletedBy     string `json:"deleted_by"`
	CreatedAt     string `json:"created_at"`
	DeletedAt     string `json:"deleted_at"`
}

type NewCompany struct {
	Id           string `json:"id"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	LegalName    string `json:"legal_name"`
	LegalAddress string `json:"legal_address"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
	TIN          string `json:"tin"`
	IBT          string `json:"ibt"`
}

type NullCompany struct {
	Id           string
	Name         string
	LegalName    sql.NullString
	Email        sql.NullString
	LegalAddress sql.NullString
	Country      sql.NullString
	ZipCode      sql.NullString
	TaxPayerId   sql.NullString
	IBT          sql.NullString
	Owner        *common.ShortUser
	Size         *corporate_service.ShortCompanySize
	CreatedBy    *common.ShortUser
	SizeId       *corporate_service.ShortCompanySize
	CreatedAt    string
}

type CompanyDefaults struct {
	Shop         *common.ShopCreatedModel
	Cashbox      *common.CashboxCreatedModel
	Cheques      []*common.ChequeCopyRequest
	PaymentTypes []*common.CommonPaymentTypes
}
