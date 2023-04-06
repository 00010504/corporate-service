package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/events/topics"
	"github.com/pkg/errors"
)

func (c *corporateService) CreateCompany(ctx context.Context, req *corporate_service.CreateCompanyRequest) (*common.ResponseID, error) {

	tr, err := c.strg.WithTransaction()
	if err != nil {
		return nil, errors.Wrap(err, "error while creating WithTransaction")
	}

	defer func() {
		if err != nil {
			_ = tr.Rollback()
		} else {
			_ = tr.Commit()
		}
	}()

	res, defaults, err := tr.Company().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating company")
	}

	if err := c.kafka.Push("v1.corporate_service.company.created.success", common.CompanyCreatedModel{
		Name:         req.Name,
		Id:           res.Id,
		CreatedBy:    req.UserId,
		Shop:         defaults.Shop,
		Cashbox:      defaults.Cashbox,
		Cheques:      defaults.Cheques,
		PaymentTypes: defaults.PaymentTypes,
	}); err != nil {
		return nil, errors.Wrap(err, "error while creating company copy")
	}

	return res, nil
}

func (c *corporateService) UpdateCompany(ctx context.Context, req *corporate_service.UpdateCompanyRequest) (*common.ResponseID, error) {

	res, err := c.strg.Company().Update(req)
	if err != nil {
		return nil, err
	}

	err = c.kafka.Push(topics.CompanyUpdateTopic, req)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *corporateService) GetCompanyById(ctx context.Context, req *common.RequestID) (*corporate_service.Company, error) {

	return c.strg.Company().GetById(req)
}

func (c *corporateService) GetAllCompanies(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllCompaniesResponse, error) {

	company, err := c.strg.Company().GetAll(req)
	if err != nil {
		return nil, err
	}

	return &corporate_service.GetAllCompaniesResponse{
			Data:  company.Data,
			Total: company.Total,
		},
		nil
}

func (c *corporateService) DeleteCompany(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	id, err := c.strg.Company().Delete(req)
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil

}
