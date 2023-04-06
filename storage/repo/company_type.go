package repo

import (
	"genproto/common"

	"genproto/corporate_service"
)

type CompanyTypeI interface {
	Create(entity *corporate_service.CreateCompanyTypeRequest) (*common.ResponseID, error)
	GetById(entity *common.RequestID) (*corporate_service.CompanyType, error)
	Update(entity *corporate_service.UpdateCompanyTypeRequest) (string, error)
	Delete(entity *common.RequestID) (string, error)
	GetAll(entity *common.SearchRequest) (*corporate_service.GetAllCompanyTypeResponse, error)
}
