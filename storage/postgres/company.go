package postgres

import (
	"database/sql"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/models"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type company struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewCompanyRepo(log logger.Logger, db *sqlx.DB) repo.CompanyI {
	return &company{
		db:  db,
		log: log,
	}
}

func (p *company) Create(entity *corporate_service.CreateCompanyRequest) (*common.ResponseID, *models.CompanyDefaults, error) {

	var (
		id            = uuid.New().String()
		userReq       = &common.Request{CompanyId: id, UserId: entity.UserId}
		shop          = common.ShopCreatedModel{Request: userReq}
		cashbox       = common.CashboxCreatedModel{Request: userReq}
		cheques       = make([]*common.ChequeCopyRequest, 0)
		chequeFields  = make(map[string]*common.FieldsValuesCopy)
		chequeIds     []string
		payment_types = make([]*common.CommonPaymentTypes, 0)
	)

	query := `
		INSERT INTO
			"company" 
		(
			"id",
			"name", 
			"type_id",
			"size_id",
			"owner",
			"created_by"
		)
		VALUES ($1, $2, $3, $4, $5, $5);
		`

	_, err := p.db.Exec(
		query,
		id,
		entity.Name,
		entity.TypeId,
		entity.SizeId,
		entity.UserId,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while inserting company")
	}

	// get company default shop

	query = `
		SELECT
			id,
			title
		FROM "shop"
		WHERE company_id = $1 AND deleted_at = 0
	`
	if err = p.db.QueryRow(query, id).Scan(&shop.Id, &shop.Name); err != nil {
		return nil, nil, errors.Wrap(err, "error while getting company shop")
	}

	// get company default cashbox

	query = `
		SELECT
			id,
			title,
			cheque_id,
			shop_id
		FROM "cashbox"
		WHERE company_id = $1 AND deleted_at = 0
	`
	if err := p.db.QueryRow(query, id).Scan(&cashbox.Id, &cashbox.Title, &cashbox.ChequeId, &cashbox.ShopId); err != nil {
		return nil, nil, errors.Wrap(err, "error while getting company cashbox")
	}

	// get company default cheques

	query = `
		SELECT 
			c.id,
			c.name,
			c.message,
			c.company_id
		FROM "cheque" c
		WHERE c.company_id = $1 AND c.deleted_at = 0;
	`

	rows, err := p.db.Query(query, id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while getting company cheques")
	}

	defer rows.Close()

	for rows.Next() {

		var cheque common.ChequeCopyRequest

		if err := rows.Scan(&cheque.Id, &cheque.Name, &cheque.Message, &cheque.CompanyId); err != nil {
			return nil, nil, err
		}

		cheque.Request = userReq
		cheque.ChequeFields = []*common.FieldsValuesCopy{}

		chequeIds = append(chequeIds, cheque.Id)
		cheques = append(cheques, &cheque)
	}

	// get cheque fields

	query = `
		SELECT field_id, cheque_id, position, is_added
		FROM "cheque_field"
		WHERE cheque_id = ANY($1);
	`

	rows, err = p.db.Query(query, pq.Array(chequeIds))
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while getting cheque_fields of cheques")
	}

	defer rows.Close()

	for rows.Next() {

		var (
			field_id  string
			cheque_id string
			position  int32
			is_added  bool
		)

		if err := rows.Scan(&field_id, &cheque_id, &position, &is_added); err != nil {
			return nil, nil, err
		}

		chequeFields[cheque_id] = &common.FieldsValuesCopy{
			FieldId:  field_id,
			Position: position,
			IsAdded:  is_added,
		}
	}

	for _, cheque := range cheques {
		cheque.ChequeFields = append(cheque.ChequeFields, chequeFields[cheque.Id])
	}

	// get payment_types

	query = `
		SELECT
			id,
			name,
			logo
		FROM "payment_type"
		WHERE company_id = $1 AND deleted_at = 0
	`
	rows, err = p.db.Query(query, id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while getting payment_types")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			p    = common.CommonPaymentTypes{Request: userReq}
			logo sql.NullString
		)

		if err := rows.Scan(&p.Id, &p.Name, &logo); err != nil {
			return nil, nil, err
		}

		p.Logo = logo.String

		payment_types = append(payment_types, &p)
	}
	return &common.ResponseID{Id: id}, &models.CompanyDefaults{Shop: &shop, Cashbox: &cashbox, Cheques: cheques, PaymentTypes: payment_types}, err
}

// ////GetById
func (p *company) GetById(req *common.RequestID) (*corporate_service.Company, error) {

	var (
		company     models.NullCompany
		user, owner models.NullShortUser
	)

	query := `
		SELECT 
			id,
			name, 
			owner,
			email,
			legal_name,
			legal_adress,
			country,
			zip_code,
			tax_payer_id,
			"ibt",
			created_at,
			created_by
		FROM 
			"company"
		WHERE 
			id = $1
	`

	err := p.db.QueryRow(query, req.Id).Scan(
		&company.Id,
		&company.Name,
		&owner.ID,
		&company.Email,
		&company.LegalName,
		&company.LegalAddress,
		&company.Country,
		&company.ZipCode,
		&company.TaxPayerId,
		&company.IBT,
		&company.CreatedAt,
		&user.ID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting company")
	}

	if user.ID.Valid {
		company.CreatedBy = &common.ShortUser{
			Id: user.ID.String,
		}
	}

	if owner.ID.Valid {
		company.Owner = &common.ShortUser{
			Id: user.ID.String,
		}
	}

	return &corporate_service.Company{
		Id:           company.Id,
		Name:         company.Name,
		LegalName:    company.LegalName.String,
		CreatedAt:    company.CreatedAt,
		Email:        company.Email.String,
		LegalAddress: company.LegalAddress.String,
		Country:      company.Country.String,
		ZipCode:      company.ZipCode.String,
		TaxPayerId:   company.TaxPayerId.String,
		Ibt:          company.IBT.String,
		Owner:        company.Owner,
		Size:         company.Size,
		CreatedBy:    company.CreatedBy,
		SizeId:       company.SizeId,
	}, nil
}

func (p *company) Update(entity *corporate_service.UpdateCompanyRequest) (*common.ResponseID, error) {
	var query = `
		UPDATE	
			"company"
		SET
			name = $2,
			email = $3,
			legal_name = $4,
			legal_adress = $5,
			country = $6,
			zip_code = $7,
			tax_payer_id = $8,
			"ibt"=$9
		WHERE
			id = $1 AND deleted_at = 0
	`
	res, err := p.db.Exec(query, entity.Id, entity.Name, entity.Email, entity.LegalName, entity.LegalAddress, entity.CountryId, entity.ZipCode, entity.TaxPayerId, entity.Ibt)

	if err != nil {
		return nil, errors.Wrap(err, "error while updating company")
	}

	if i, _ := res.RowsAffected(); i == 0 {
		return nil, sql.ErrNoRows
	}

	return &common.ResponseID{Id: entity.Id}, nil
}

func (p *company) GetAll(req *common.SearchRequest) (*corporate_service.GetAllCompaniesResponse, error) {
	var (
		data         []*corporate_service.Company = make([]*corporate_service.Company, 0)
		company      corporate_service.Company
		searchFields map[string]interface{} = map[string]interface{}{
			"limit":      req.Limit,
			"offset":     req.Limit * (req.Page - 1),
			"long_name":  req.Search,
			"short_name": req.Search,
		}
	)

	filter := `
		WHERE
			deleted_at = 0
	`
	limit := ` LIMIT :limit`
	offset := " OFFSET :offset;"

	rQuery := `
		SELECT 
			id,
			name,
			owner,
			company_type_id,
			created_at,
			created_by
		FROM 
			"company"
	`

	if req.GetSearch() != "" {
		filter += `
		AND
		(
			long_name ILIKE '%' || :long_name || '%'
			OR
			short_name ILIKE '%' || :short_name || '%'
		)
	`
	}

	rQuery += filter + limit + offset

	rows, err := p.db.NamedQuery(rQuery, searchFields)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting company pagination")
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&company.Id,
			&company.Name,
			&company.Owner,
			&company,
			&company.CreatedAt,
			&company.CreatedBy,
		)

		if err != nil {
			return nil, errors.Wrap(err, "error while getting company pagination rows.Scan")
		}

		data = append(data, &corporate_service.Company{
			Id:        company.Id,
			Name:      company.Name,
			Owner:     company.Owner,
			CreatedAt: company.CreatedAt,
			CreatedBy: company.CreatedBy,
		})
	}

	var total int64
	countQuery := `
		SELECT 
			COUNT(1)
		FROM 
			"company"
	`
	countQuery += filter

	cRows, err := p.db.NamedQuery(countQuery, searchFields)
	if err != nil {
		return nil, errors.Wrap(err, "error while scanning count")
	}

	defer cRows.Close()

	for cRows.Next() {
		err = cRows.Scan(&total)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning count")
		}
	}

	return &corporate_service.GetAllCompaniesResponse{
		Data:  data,
		Total: total,
	}, nil
}

func (p *company) Delete(entity *common.RequestID) (string, error) {

	var query = `UPDATE company 
						SET deleted_at=extract(epoch from now())::bigint 
								WHERE id = $1 and deleted_at = 0`

	res, err := p.db.Exec(query, entity.Id)
	if err != nil {
		return "", errors.Wrap(err, "error while deleting company")
	}

	i, err := res.RowsAffected()
	if err != nil {
		return "", errors.Wrap(err, "error while deleting company")
	}

	if i == 0 {
		return "", errors.New("company not found")
	}

	return entity.Id, nil
}
