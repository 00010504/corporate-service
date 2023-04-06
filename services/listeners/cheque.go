package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/pkg/errors"
)

func (c *corporateService) CreateCheque(ctx context.Context, req *corporate_service.CreateChequeRequest) (*common.ResponseID, error) {

	var (
		chequeFields []*common.FieldsValuesCopy
	)

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

	res, err := c.strg.Cheque().Create(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating cheque")
	}

	for _, field := range req.FieldIds {
		chequeFields = append(chequeFields, &common.FieldsValuesCopy{
			FieldId:  field.FieldId,
			Position: field.Position,
			IsAdded:  field.IsAdded,
		})
	}

	err = c.kafka.Push("v1.corporate_service.cheque.created.success", common.ChequeCopyRequest{
		Id:        res.Id,
		CompanyId: req.Request.CompanyId,
		Name:      req.Name,
		Message:   req.Message,
		ChequeLogo: &common.ChequeLogoCopyRequest{
			Image:    req.Logo.Image,
			ChequeId: res.Id,
			Left:     int32(req.Logo.Left),
			Right:    int32(req.Logo.Right),
			Top:      int32(req.Logo.Top),
			Bottom:   int32(req.Logo.Bottom),
		},
		ChequeFields: chequeFields,
		Request:      req.Request,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *corporateService) GetAllCheques(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllChequesResponse, error) {
	return c.strg.Cheque().GetAll(req)
}

func (c *corporateService) GetCheque(ctx context.Context, req *common.RequestID) (*corporate_service.Cheque, error) {
	return c.strg.Cheque().Get(req)
}

func (c *corporateService) GetAllRecieptBlock(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllRecieptBlockResponse, error) {
	return c.strg.Cheque().GetAllRecieptBlock(req)
}

func (c *corporateService) DeleteCheque(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {
	return c.strg.Cheque().Delete(ctx, req)
}

func (c *corporateService) UpdateCheque(ctx context.Context, req *corporate_service.UpdateChequeRequest) (*common.ResponseID, error) {

	var (
		chequeFields []*common.FieldsValuesCopy
	)

	res, err := c.strg.Cheque().Update(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, field := range req.FieldIds {
		chequeFields = append(chequeFields, &common.FieldsValuesCopy{
			FieldId:  field.FieldId,
			Position: field.Position,
		})
	}

	err = c.kafka.Push("v1.corporate_service.cheque.created.success", common.ChequeCopyRequest{
		Id:        res.Id,
		CompanyId: req.Request.CompanyId,
		Name:      req.Name,
		Message:   req.Message,
		ChequeLogo: &common.ChequeLogoCopyRequest{
			Image:    req.Logo.Image,
			ChequeId: res.Id,
			Left:     int32(req.Logo.Left),
			Right:    int32(req.Logo.Right),
			Top:      int32(req.Logo.Top),
			Bottom:   int32(req.Logo.Bottom),
		},
		ChequeFields: chequeFields,
		Request:      req.Request,
	})

	if err != nil {
		return nil, err
	}
	return res, nil
}
