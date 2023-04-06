package repo

import (
	"genproto/common"

	"genproto/corporate_service"
)

type ShopI interface {
	Create(entity *corporate_service.CreateShopRequest) (string, error)
	GetById(entity *common.RequestID) (*corporate_service.Shop, error)
	Update(entity *corporate_service.UpdateShopRequest) (string, error)
	Delete(entity *common.RequestID) (string, error)
	GetAll(req *common.SearchRequest) (*corporate_service.GetAllShopsResponse, error)
}
