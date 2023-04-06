package repo

import (
	"genproto/common"

	"genproto/corporate_service"
)

type PaymentTypeI interface {
	Create(entity *corporate_service.CreatePaymentTypeRequest) (string, error)
	Update(entity *corporate_service.UpdatePaymentTypeRequest) (string, error)
	GetById(id string) (*corporate_service.PaymentType, error)
	GetAll(req *common.SearchRequest) (*corporate_service.GetPaymentTypesResponse, error)
	Delete(entity *common.RequestID) (string, error)
}
