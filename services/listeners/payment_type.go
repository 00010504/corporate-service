package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/pkg/errors"
)

func (c *corporateService) CreatePaymentType(ctx context.Context, req *corporate_service.CreatePaymentTypeRequest) (*common.ResponseID, error) {

	id, err := c.strg.PaymentType().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating payment type")
	}
	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) UpdatePaymentType(ctx context.Context, req *corporate_service.UpdatePaymentTypeRequest) (*common.ResponseID, error) {

	id, err := c.strg.PaymentType().Update(req)

	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) GetPaymentTypes(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetPaymentTypesResponse, error) {

	return c.strg.PaymentType().GetAll(req)

}

func (c *corporateService) GetPaymentTypeById(ctx context.Context, req *common.RequestID) (*corporate_service.PaymentType, error) {
	paymentType, err := c.strg.PaymentType().GetById(req.Id)
	if err != nil {
		return nil, err
	}
	return &corporate_service.PaymentType{
		Id:        paymentType.Id,
		Name:      paymentType.Name,
		CreatedAt: paymentType.CreatedAt,
	}, nil
}

func (c *corporateService) GetAllPaymentTypes(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetPaymentTypesResponse, error) {

	return c.strg.PaymentType().GetAll(req)
}

func (c *corporateService) DeletePaymentType(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	id, err := c.strg.PaymentType().Delete(req)
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil

}
