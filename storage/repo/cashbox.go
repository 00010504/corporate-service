package repo

import (
	"genproto/common"

	"genproto/corporate_service"
)

type CashboxI interface {
	Create(entity *corporate_service.CreateCashboxRequest) (string, error)
	GetById(req *common.RequestID) (*corporate_service.Cashbox, error)
	Update(entity *corporate_service.UpdateCashboxRequest) (*common.ResponseID, error)
	GetAll(req *common.ShopSearchRequest) (*corporate_service.GetAllCashboxesResponse, error)
	Delete(entity *common.RequestID) (string, error)
}
