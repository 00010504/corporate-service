package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/pkg/errors"
)

func (c *corporateService) CreateCompanyType(ctx context.Context, req *corporate_service.CreateCompanyTypeRequest) (*common.ResponseID, error) {

	res, err := c.strg.CompanyType().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating company type")
	}

	return res, nil
}

func (c *corporateService) UpdateCompanyType(ctx context.Context, req *corporate_service.UpdateCompanyTypeRequest) (*common.ResponseID, error) {

	id, err := c.strg.CompanyType().Update(req)

	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) GetCompanyTypeById(ctx context.Context, req *common.RequestID) (*corporate_service.CompanyType, error) {

	companyType, err := c.strg.CompanyType().GetById(req)
	if err != nil {
		return nil, err
	}
	return companyType, nil
}

func (c *corporateService) DeleteCompanyType(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	id, err := c.strg.CompanyType().Delete(req)
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil

}

func (s *corporateService) GetAllCompanyTypes(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllCompanyTypeResponse, error) {
	return s.strg.CompanyType().GetAll(req)
}
