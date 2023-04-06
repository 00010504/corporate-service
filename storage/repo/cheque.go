package repo

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"
)

type ChequeI interface {
	Create(entity *corporate_service.CreateChequeRequest) (*common.ResponseID, error)
	GetAll(entity *common.SearchRequest) (*corporate_service.GetAllChequesResponse, error)
	Get(entity *common.RequestID) (*corporate_service.Cheque, error)
	GetAllRecieptBlock(req *common.SearchRequest) (*corporate_service.GetAllRecieptBlockResponse, error)
	Delete(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)
	Update(ctx context.Context, req *corporate_service.UpdateChequeRequest) (*common.ResponseID, error)
}
