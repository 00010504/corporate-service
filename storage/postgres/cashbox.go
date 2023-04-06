package postgres

import (
	"context"
	"database/sql"
	"genproto/common"
	"strings"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/models"
	"github.com/Invan2/invan_corporate_service/pkg/helper"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type cashboxRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewCashbox(log logger.Logger, db *sqlx.DB) repo.CashboxI {
	return &cashboxRepo{
		db:  db,
		log: log,
	}
}

func (r *cashboxRepo) Create(entity *corporate_service.CreateCashboxRequest) (string, error) {

	r.log.Info("create cashbox", logger.Any("cashbox", entity))

	var cashboxId = uuid.New().String()

	tr, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return "", errors.Wrap(err, "error while create transaction")
	}

	defer func() {
		if err != nil {
			tr.Rollback()
		} else {
			tr.Commit()
		}
	}()

	query := `
		INSERT INTO
			"cashbox" 
		(
			id, 
			shop_id, 
			cheque_id,
			title, 
			created_by,
			company_id
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
	`

	_, err = tr.Exec(
		query,
		cashboxId,
		helper.NullString(entity.ShopId),
		entity.ReceiptId,
		entity.Title,
		entity.Request.UserId,
		entity.Request.CompanyId,
	)
	if err != nil {
		return "", errors.Wrap(err, "error while creating cashbox")
	}

	if len(entity.PaymentTypeIds) > 0 {

		var (
			values = []interface{}{}
			stmt   *sql.Stmt
		)

		query = `
		INSERT INTO
			"cashbox_payment"
		(
			id,
			cashbox_id,
			payment_type_id,
			created_by
		)
		VALUES 
	`

		for _, paymentTypeId := range entity.PaymentTypeIds {
			query += "(?, ?, ?, ?),"
			values = append(
				values,
				uuid.New().String(),
				cashboxId,
				paymentTypeId,
				entity.Request.UserId,
			)
		}

		query = strings.TrimSuffix(query, ",")
		query = helper.ReplaceSQL(query, "?")

		stmt, err := tr.Prepare(query)
		if err != nil {
			return "", errors.Wrap(err, "error while creating cashbox_payment Prepare")
		}

		defer stmt.Close()

		_, err = stmt.Exec(values...)
		if err != nil {
			return "", errors.Wrap(err, "error while creating cashbox_payment Exec")
		}
		stmt.Close()
	}

	return cashboxId, nil
}

func (r *cashboxRepo) GetById(req *common.RequestID) (*corporate_service.Cashbox, error) {

	var (
		cashbox = corporate_service.Cashbox{
			Cheques:        make([]*corporate_service.ShortCheque, 0),
			PaymentTypes:   make([]*corporate_service.ShortPaymentType, 0),
			PaymentTypeIds: make([]string, 0),
			CreatedBy:      &common.ShortUser{},
			Shops:          make([]*common.ShortShop, 0),
		}
		shopId   sql.NullString
		chequeID sql.NullString
	)

	query := `
		SELECT 
			c.id, 
			c.shop_id, 
			c.title, 
			c.cheque_id, 
			c.created_at, 
			c.created_by
		FROM
			"cashbox" c
		WHERE
			c.id = $1 AND c.company_id=$2 AND c.deleted_at = 0
		`

	err := r.db.QueryRow(query, req.Id, req.Request.CompanyId).Scan(
		&cashbox.Id,
		&shopId,
		&cashbox.Title,
		&chequeID,
		&cashbox.CreatedAt,
		&cashbox.CreatedBy.Id,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting cashbox")
	}

	cashbox.ChequeId = chequeID.String

	query = `
		SELECT 
			"id",
			"name"
		FROM "cheque"
		WHERE company_id=$1 and deleted_at=0
	`

	rows, err := r.db.Query(query, req.Request.CompanyId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			cheque corporate_service.ShortCheque
		)

		err := rows.Scan(&cheque.Id, &cheque.Name)
		if err != nil {
			return nil, err
		}

		cheque.IsAdded = (cheque.Id == chequeID.String)

		cashbox.Cheques = append(cashbox.Cheques, &cheque)
	}

	rows.Close()

	query = `
		SELECT
			pt."id",
			pt."name",
			cp.id
		FROM payment_type pt
		LEFT JOIN "cashbox_payment" cp ON cp.payment_type_id = pt.id AND cp.cashbox_id=$2 AND cp.deleted_at=0
		WHERE pt.company_id=$1 and pt.deleted_at=0
	`

	rows, err = r.db.Query(query, req.Request.CompanyId, cashbox.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			paymentType corporate_service.ShortPaymentType
			id          sql.NullString
		)

		err := rows.Scan(&paymentType.Id, &paymentType.Name, &id)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			cashbox.PaymentTypeIds = append(cashbox.PaymentTypeIds, paymentType.Id)
		}

		paymentType.IsAdded = id.Valid

		cashbox.PaymentTypes = append(cashbox.PaymentTypes, &paymentType)

	}

	rows.Close()

	query = `
		SELECT
			"id",
			"title"
		FROM shop
		WHERE company_id=$1 AND deleted_at=0
	`

	rows, err = r.db.Query(query, req.Request.CompanyId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var (
			shop common.ShortShop
		)

		err := rows.Scan(&shop.Id, &shop.Name)
		if err != nil {
			return nil, err
		}

		shop.IsAdded = shop.Id == shopId.String
		if shop.IsAdded {
			cashbox.Shop = &shop
		}

		cashbox.Shops = append(cashbox.Shops, &shop)
	}

	return &cashbox, nil
}

func (r *cashboxRepo) Update(entity *corporate_service.UpdateCashboxRequest) (*common.ResponseID, error) {

	r.log.Info("update cahsbox", logger.Any("entitry", entity))

	tr, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {

		if err != nil {
			_ = tr.Rollback()
		} else {
			_ = tr.Commit()
		}

	}()

	query := `
		UPDATE "cashbox"
			SET 
				shop_id = $2, 
				title = $3,
				cheque_id = $4
				
		WHERE id = $1 AND company_id=$5 AND deleted_at=0`

	res, err := tr.Exec(query, entity.Id, helper.NullString(entity.ShopId), entity.Title, entity.ChequeId, entity.Request.CompanyId)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating cashbox")
	}

	if i, _ := res.RowsAffected(); i == 0 {
		err = sql.ErrNoRows
		return nil, err
	}

	query = `
		UPDATE "cashbox_payment" SET deleted_at=extract(epoch from now())::bigint  WHERE cashbox_id=$1 AND deleted_at=0
	
	`
	_, err = tr.Exec(query, entity.Id)
	if err != nil {
		return nil, err
	}

	if len(entity.PaymentTypeIds) > 0 {

		var (
			values = []interface{}{}
			stmt   *sql.Stmt
		)

		query = `
		INSERT INTO
			"cashbox_payment"
		(
			id,
			cashbox_id,
			payment_type_id,
			created_by
		)
		VALUES 
	`

		for _, paymentTypeId := range entity.PaymentTypeIds {
			query += "(?, ?, ?, ?),"
			values = append(
				values,
				uuid.New().String(),
				entity.Id,
				paymentTypeId,
				entity.Request.UserId,
			)
		}

		query = strings.TrimSuffix(query, ",")
		query = helper.ReplaceSQL(query, "?")

		stmt, err = tr.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "error while creating cashbox_payment Prepare")
		}

		defer stmt.Close()

		_, err = stmt.Exec(values...)
		if err != nil {
			return nil, errors.Wrap(err, "error while creating cashbox_payment Exec")
		}
	}

	return &common.ResponseID{Id: entity.Id}, nil
}

func (r *cashboxRepo) GetAll(req *common.ShopSearchRequest) (*corporate_service.GetAllCashboxesResponse, error) {

	var (
		res corporate_service.GetAllCashboxesResponse = corporate_service.GetAllCashboxesResponse{
			Data: make([]*corporate_service.ShortCashbox, 0),
		}

		searchFields map[string]interface{} = map[string]interface{}{
			"limit":      req.Limit,
			"offset":     req.Limit * (req.Page - 1),
			"search":     req.Search,
			"company_id": req.Request.CompanyId,
			"shop_id":    req.ShopId,
		}
	)

	filter := `
		WHERE
			c.deleted_at = 0 AND c.company_id = :company_id
	`

	rQuery := `
		SELECT 
			c.id,
			c.title,
			c.created_at,
			sh.id,
			sh.title,
			c.created_by,
			u.first_name,
			u.last_name
		FROM 
			"cashbox" c
		LEFT JOIN "shop" sh ON sh.id = c.shop_id AND sh.deleted_at = 0
		LEFT JOIN "user" u ON u.id = c.created_by AND u.deleted_at = 0
	`

	if req.GetSearch() != "" {
		filter += `
		AND
		(
			c.title ILIKE '%' || :search || '%' OR 
			sh.title ILIKE '%' || :search || '%' 
			
		)
	`
	}

	if req.GetShopId() != "" {
		filter += `
		AND
		sh.id =:shop_id
	`
	}

	rQuery += filter + `ORDER BY c.created_at DESC LIMIT :limit OFFSET :offset`

	rows, err := r.db.NamedQuery(rQuery, searchFields)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting cashbox")
	}

	defer rows.Close()

	for rows.Next() {

		var (
			cashbox corporate_service.ShortCashbox
			shop    models.NullShortShop
			user    models.NullShortUser
		)

		err = rows.Scan(
			&cashbox.Id,
			&cashbox.Name,
			&cashbox.CreatedAt,
			&shop.ID,
			&shop.Name,
			&user.ID,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting cashbox pagination rows.Scan")
		}

		if shop.ID.Valid {
			cashbox.Shop = &common.ShortShop{
				Id:   shop.ID.String,
				Name: shop.Name.String,
			}
		}

		if user.ID.Valid {
			cashbox.CreatedBy = &common.ShortUser{
				Id:        user.ID.String,
				FirstName: user.FirstName.String,
				LastName:  user.LastName.String,
			}
		}

		res.Data = append(res.Data, &cashbox)
	}

	countQuery := `
		SELECT 
			count(*)
		FROM 
			"cashbox" c
		LEFT JOIN "shop" sh ON sh.id = c.shop_id AND sh.deleted_at = 0
		LEFT JOIN "user" u ON u.id = c.created_by AND u.deleted_at = 0
	`
	countQuery += filter

	stmt, err := r.db.PrepareNamed(countQuery)
	if err != nil {
		return nil, errors.Wrap(err, "error while scanning count")
	}

	defer stmt.Close()

	if err := stmt.QueryRow(searchFields).Scan(&res.Total); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *cashboxRepo) Delete(entity *common.RequestID) (string, error) {

	var query = `
		UPDATE
			"cashbox" 
		SET
			deleted_at=extract(epoch from now())::bigint 
		WHERE
			deleted_at = 0 AND id = $1
	`

	res, err := p.db.Exec(query, entity.Id)
	if err != nil {
		return "", errors.Wrap(err, "error while deleting cashbox")
	}

	i, err := res.RowsAffected()
	if err != nil {
		return "", errors.Wrap(err, "error while deleting cashbox")
	}

	if i == 0 {
		return "", errors.New("cashbox not found")
	}

	return entity.Id, nil
}
