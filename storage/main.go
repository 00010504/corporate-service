package storage

import (
	"context"
	"database/sql"

	"github.com/Invan2/invan_corporate_service/config"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/postgres"
	"github.com/Invan2/invan_corporate_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type repos struct {
	companyRepo     repo.CompanyI
	userRepo        repo.UserI
	paymentTypeRepo repo.PaymentTypeI
	companyTypeRepo repo.CompanyTypeI
	companySizeRepo repo.CompanySizeI
	shopRepo        repo.ShopI
	cashboxRepo     repo.CashboxI
	chequeRepo      repo.ChequeI
}

type repoIs interface {
	Company() repo.CompanyI
	User() repo.UserI
	PaymentType() repo.PaymentTypeI
	CompanyType() repo.CompanyTypeI
	CompanySize() repo.CompanySizeI
	Shop() repo.ShopI
	Cashbox() repo.CashboxI
	Cheque() repo.ChequeI
}

type storage struct {
	db  *sqlx.DB
	log logger.Logger
	repos
}

type storageTr struct {
	tr *sqlx.Tx
	repos
}

type StorageTrI interface {
	Commit() error
	Rollback() error
	repoIs
}

type StoragePgI interface {
	WithTransaction() (StorageTrI, error)
	repoIs
}

func NewStoragePg(log logger.Logger, db *sqlx.DB, cfg config.Config) StoragePgI {

	return &storage{
		db:  db,
		log: log,
		repos: repos{

			companyRepo:     postgres.NewCompanyRepo(log, db),
			userRepo:        postgres.NewUserRepo(log, db),
			shopRepo:        postgres.NewShopRepo(log, db),
			paymentTypeRepo: postgres.NewPaymentType(log, db),
			cashboxRepo:     postgres.NewCashbox(log, db),
			companyTypeRepo: postgres.NewCompanyType(log, db),
			companySizeRepo: postgres.NewCompanySizeRepo(log, db),
			chequeRepo:      postgres.NewChequeRepo(db, log, cfg),
		},
	}
}

func (s *storage) WithTransaction() (StorageTrI, error) {

	tr, err := s.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	return &storageTr{
		tr:    tr,
		repos: s.repos,
	}, nil
}

func (s *storageTr) Commit() error {
	return s.tr.Commit()
}

func (s *storageTr) Rollback() error {
	return s.tr.Rollback()
}

func (r *repos) Company() repo.CompanyI {
	return r.companyRepo
}

func (r *repos) CompanySize() repo.CompanySizeI {
	return r.companySizeRepo
}

func (r *repos) CompanyType() repo.CompanyTypeI {
	return r.companyTypeRepo
}

func (r *repos) User() repo.UserI {
	return r.userRepo
}

func (r *repos) PaymentType() repo.PaymentTypeI {
	return r.paymentTypeRepo
}

func (r *repos) Shop() repo.ShopI {
	return r.shopRepo
}

func (r *repos) Cashbox() repo.CashboxI {
	return r.cashboxRepo
}

func (r *repos) Cheque() repo.ChequeI {
	return r.chequeRepo
}
