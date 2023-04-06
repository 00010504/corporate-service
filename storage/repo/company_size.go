package repo

import (
	"genproto/common"

	"genproto/corporate_service"
)

type CompanySizeI interface {
	Create(entity *corporate_service.CreateCompanySizeRequest) (*common.ResponseID, error)
	GetAll(entity *common.SearchRequest) (*corporate_service.GetAllCompanySizeResponse, error)
}
