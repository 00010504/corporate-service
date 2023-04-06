package repo

import (
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/models"
)

type CompanyI interface {
	Create(entity *corporate_service.CreateCompanyRequest) (*common.ResponseID, *models.CompanyDefaults, error)
	GetById(req *common.RequestID) (*corporate_service.Company, error)
	Update(req *corporate_service.UpdateCompanyRequest) (*common.ResponseID, error)
	GetAll(req *common.SearchRequest) (*corporate_service.GetAllCompaniesResponse, error)
	Delete(entity *common.RequestID) (string, error)
}
