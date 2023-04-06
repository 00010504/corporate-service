package postgres

import (
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type paymentTypeRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewPaymentType(log logger.Logger, db *sqlx.DB) repo.PaymentTypeI {
	return &paymentTypeRepo{
		db:  db,
		log: log,
	}
}

// /////Create
func (p *paymentTypeRepo) Create(entity *corporate_service.CreatePaymentTypeRequest) (string, error) {
	var query = `
		INSERT INTO
				"payment" 
			(
					id,
					name
			)
			VALUES 
			(
				$1,
				$2

			);
		`
	entity.Id = uuid.New().String()

	_, err := p.db.Query(
		query,
		entity.Id,
		entity.Name,
	)
	if err != nil {
		return "", errors.Wrap(err, "error while inserting payment type")
	}
	return entity.Id, err
}

// ////GetById
func (p *paymentTypeRepo) GetById(id string) (*corporate_service.PaymentType, error) {
	var paymentType corporate_service.PaymentType

	var query = `
		SELECT 
			id,
			name, 
			created_at,
			deleted_at
		FROM 
			"payment"
		WHERE 
			id = $1
	`
	err := p.db.QueryRow(query, id).Scan(&paymentType.Id, &paymentType.Name, &paymentType.CreatedAt, &paymentType.CreatedBy)

	if err != nil {
		return nil, errors.Wrap(err, "error while getting payment type")
	}

	return &paymentType, nil
}

/////get all

func (p *paymentTypeRepo) GetAll(req *common.SearchRequest) (*corporate_service.GetPaymentTypesResponse, error) {
	var (
		searchFields map[string]interface{} = map[string]interface{}{
			"limit":      req.Limit,
			"offset":     req.Limit * (req.Page - 1),
			"search":     req.Search,
			"company_id": req.Request.CompanyId,
		}
	)
	res := corporate_service.GetPaymentTypesResponse{
		Data: make([]*corporate_service.PaymentType, 0),
	}

	filter := `
		WHERE
			deleted_at = 0 and company_id = :company_id
	`
	limit := ` LIMIT :limit`
	offset := " OFFSET :offset;"

	rQuery := `
		SELECT 
			pt.id,
			pt.name,
			pt.created_at
		FROM 
			"payment_type" AS pt
	`

	if req.GetSearch() != "" {
		filter += `
		AND
		(
			pt.name ILIKE '%' || :search || '%'
			
		)
	`
	}

	rQuery += filter + limit + offset

	rows, err := p.db.NamedQuery(rQuery, searchFields)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting payment type pagination")
	}

	defer rows.Close()
	for rows.Next() {

		var paymentType corporate_service.PaymentType

		err = rows.Scan(
			&paymentType.Id,
			&paymentType.Name,
			&paymentType.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting payment type  pagination rows.Scan")
		}

		res.Data = append(res.Data, &corporate_service.PaymentType{
			Id:        paymentType.Id,
			Name:      paymentType.Name,
			CreatedAt: paymentType.CreatedAt,
		})
	}

	countQuery := `
		SELECT 
			COUNT(*) AS total
		FROM 
			"payment_type" AS pt
	`
	countQuery += filter

	stmt, err := p.db.PrepareNamed(countQuery)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.QueryRow(searchFields).Scan(&res.Total); err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *paymentTypeRepo) Update(entity *corporate_service.UpdatePaymentTypeRequest) (string, error) {
	var query = `
		UPDATE	
			"payment"
		SET
			name = $2
		WHERE
			id = $1
		RETURNING 
			id;
	`
	err := p.db.QueryRow(query, entity.Id, entity.Name).Scan(&entity.Id)

	if err != nil {
		return "", errors.Wrap(err, "error while updating payment type")
	}

	return entity.Id, nil
}

func (p *paymentTypeRepo) Delete(entity *common.RequestID) (string, error) {

	var query = `UPDATE payment 
						SET deleted_at=extract(epoch from now())::bigint 
								WHERE id = $1 and deleted_at = 0 and company_id = $2`

	res, err := p.db.Exec(query, entity.Id, entity.Request.CompanyId)
	if err != nil {
		return "", errors.Wrap(err, "error while deleting payment type")
	}

	i, err := res.RowsAffected()
	if err != nil {
		return "", errors.Wrap(err, "error while deleting payment type")
	}

	if i == 0 {
		return "", errors.New("payment type not found")
	}

	return entity.Id, nil
}
