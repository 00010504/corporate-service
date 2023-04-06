package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/pkg/errors"
)

func (c *corporateService) CreateCashbox(ctx context.Context, req *corporate_service.CreateCashboxRequest) (*common.ResponseID, error) {

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

	id, err := tr.Cashbox().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating cashbox")
	}

	err = c.kafka.Push("v1.corporate_service.cashbox.created.success", common.CashboxCreatedModel{
		Id:       id,
		ShopId:   req.ShopId,
		Title:    req.Title,
		ChequeId: req.ReceiptId,
		Request:  req.Request,
	})
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) GetCashboxById(ctx context.Context, req *common.RequestID) (*corporate_service.Cashbox, error) {

	return c.strg.Cashbox().GetById(req)
}

func (c *corporateService) UpdateCashbox(ctx context.Context, req *corporate_service.UpdateCashboxRequest) (*common.ResponseID, error) {

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

	updateRes, err := tr.Cashbox().Update(req)
	if err != nil {
		return nil, err
	}

	err = c.kafka.Push("v1.corporate_service.cashbox.created.success", common.CashboxCreatedModel{
		Id:       updateRes.Id,
		ShopId:   req.ShopId,
		Title:    req.Title,
		ChequeId: req.ChequeId,
		Request:  req.Request,
	})
	if err != nil {
		return nil, err
	}

	return updateRes, nil
}

func (c *corporateService) GetAllCashboxes(ctx context.Context, req *common.ShopSearchRequest) (*corporate_service.GetAllCashboxesResponse, error) {

	return c.strg.Cashbox().GetAll(req)
}
func (c *corporateService) DeleteCashbox(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	id, err := c.strg.Cashbox().Delete(req)
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil

}
