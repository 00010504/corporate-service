package listeners

import (
	"context"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/events"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage"
)

type corporateService struct {
	log   logger.Logger
	kafka events.PubSubServer
	strg  storage.StoragePgI
}

type CorporateService interface {

	//company
	CreateCompany(ctx context.Context, req *corporate_service.CreateCompanyRequest) (*common.ResponseID, error)
	GetCompanyById(ctx context.Context, req *common.RequestID) (*corporate_service.Company, error)
	UpdateCompany(ctx context.Context, req *corporate_service.UpdateCompanyRequest) (*common.ResponseID, error)
	GetAllCompanies(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllCompaniesResponse, error)
	DeleteCompany(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)

	// company type
	CreateCompanyType(ctx context.Context, req *corporate_service.CreateCompanyTypeRequest) (*common.ResponseID, error)
	UpdateCompanyType(ctx context.Context, req *corporate_service.UpdateCompanyTypeRequest) (*common.ResponseID, error)
	GetCompanyTypeById(ctx context.Context, req *common.RequestID) (*corporate_service.CompanyType, error)
	DeleteCompanyType(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)
	GetAllCompanyTypes(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllCompanyTypeResponse, error)

	// company size
	CreateCompanySize(ctx context.Context, req *corporate_service.CreateCompanySizeRequest) (*common.ResponseID, error)
	GetAllCompanySize(ctx context.Context, entity *common.SearchRequest) (*corporate_service.GetAllCompanySizeResponse, error)

	// shop
	CreateShop(ctx context.Context, req *corporate_service.CreateShopRequest) (*common.ResponseID, error)
	GetShopById(ctx context.Context, req *common.RequestID) (*corporate_service.Shop, error)
	UpdateShop(ctx context.Context, req *corporate_service.UpdateShopRequest) (*common.ResponseID, error)
	GetAllShops(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllShopsResponse, error)
	DeleteShop(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)

	// cashbox
	CreateCashbox(ctx context.Context, req *corporate_service.CreateCashboxRequest) (*common.ResponseID, error)
	GetCashboxById(ctx context.Context, req *common.RequestID) (*corporate_service.Cashbox, error)
	UpdateCashbox(ctx context.Context, req *corporate_service.UpdateCashboxRequest) (*common.ResponseID, error)
	GetAllCashboxes(ctx context.Context, req *common.ShopSearchRequest) (*corporate_service.GetAllCashboxesResponse, error)
	DeleteCashbox(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)

	// payment type
	CreatePaymentType(ctx context.Context, req *corporate_service.CreatePaymentTypeRequest) (*common.ResponseID, error)
	UpdatePaymentType(ctx context.Context, req *corporate_service.UpdatePaymentTypeRequest) (*common.ResponseID, error)
	GetAllPaymentTypes(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetPaymentTypesResponse, error)
	DeletePaymentType(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)

	// cheque

	CreateCheque(ctx context.Context, req *corporate_service.CreateChequeRequest) (*common.ResponseID, error)
	GetAllCheques(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllChequesResponse, error)
	GetCheque(ctx context.Context, req *common.RequestID) (*corporate_service.Cheque, error)
	GetAllRecieptBlock(ctx context.Context, req *common.SearchRequest) (*corporate_service.GetAllRecieptBlockResponse, error)
	DeleteCheque(ctx context.Context, req *common.RequestID) (*common.ResponseID, error)
	UpdateCheque(ctx context.Context, req *corporate_service.UpdateChequeRequest) (*common.ResponseID, error)
}

func NewCorporateService(log logger.Logger, kafka events.PubSubServer, strg storage.StoragePgI) CorporateService {
	return &corporateService{
		log:   log,
		kafka: kafka,
		strg:  strg,
	}
}
