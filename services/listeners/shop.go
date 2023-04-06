package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/pkg/errors"
)

func (c *corporateService) CreateShop(ctx context.Context, req *corporate_service.CreateShopRequest) (*common.ResponseID, error) {

	c.log.Info("data", logger.Any("data", req))

	id, err := c.strg.Shop().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating shop")
	}

	err = c.kafka.Push("v1.corporate_service.shop.created.success", common.ShopCreatedModel{
		Name:    req.Title,
		Id:      id,
		Request: req.Request,
	})
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) GetShopById(ctx context.Context, req *common.RequestID) (*corporate_service.Shop, error) {

	shop, err := c.strg.Shop().GetById(req)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (c *corporateService) UpdateShop(ctx context.Context, req *corporate_service.UpdateShopRequest) (*common.ResponseID, error) {

	id, err := c.strg.Shop().Update(req)
	if err != nil {
		return nil, err
	}

	err = c.kafka.Push("v1.corporate_service.shop.created.success", common.ShopCreatedModel{
		Name:    req.Title,
		Id:      id,
		Request: req.Request,
	})
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}

func (c *corporateService) GetAllShops(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllShopsResponse, error) {

	res, err := c.strg.Shop().GetAll(req)
	if err != nil {
		return nil, err
	}

	c.log.Info("get all shops", logger.Any("res", res))

	return res, nil
}

func (c *corporateService) DeleteShop(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	id, err := c.strg.Shop().Delete(req)
	if err != nil {
		return nil, err
	}

	err = c.kafka.Push("v1.corporate_service.shop.deleted.success", common.RequestID{
		Id:      id,
		Request: req.Request,
	})
	if err != nil {
		return nil, err
	}

	return &common.ResponseID{Id: id}, nil
}
