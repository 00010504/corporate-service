package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"
)

func (c *corporateService) CreateCompanySize(ctx context.Context, req *corporate_service.CreateCompanySizeRequest) (*common.ResponseID, error) {
	return nil, nil
}

func (c *corporateService) GetAllCompanySize(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllCompanySizeResponse, error) {

	res, err := c.strg.CompanySize().GetAll(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
