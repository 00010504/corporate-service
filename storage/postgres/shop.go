package postgres

import (
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/models"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type shopRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewShopRepo(log logger.Logger, db *sqlx.DB) repo.ShopI {
	return &shopRepo{
		db:  db,
		log: log,
	}
}

func (p *shopRepo) Create(entity *corporate_service.CreateShopRequest) (string, error) {

	id := uuid.New().String()

	var query = `
		INSERT INTO
				shop
			(
					id,
					title, 
					company_id,
					phone_number, 
					size,
					address,
					description,
					created_by
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := p.db.Exec(
		query,
		id,
		entity.Title,
		entity.Request.CompanyId,
		entity.PhoneNumber,
		entity.Size,
		entity.Address,
		entity.Description,
		entity.Request.UserId,
	)
	if err != nil {
		return "", errors.Wrap(err, "error while inserting shop")
	}

	return id, err
}

func (p *shopRepo) GetById(entity *common.RequestID) (*corporate_service.Shop, error) {

	var (
		shop corporate_service.Shop
		user models.NullShortUser
	)

	var query = `
		SELECT 
			s.id,
			s.title, 
			s.phone_number,
			s.size,
			s.address,
			s.description,
			s.created_at,
			s.created_by,
			u.first_name,
			u.last_name
		FROM "shop" s
		LEFT JOIN "user" u ON u.id = s.created_by and u.deleted_at = 0
		WHERE 
			s.id = $1 and s.deleted_at = 0
	`
	err := p.db.QueryRow(query, entity.Id).Scan(
		&shop.Id,
		&shop.Title,
		&shop.PhoneNumber,
		&shop.Size,
		&shop.Address,
		&shop.Description,
		&shop.CreatedAt,
		&user.ID,
		&user.FirstName,
		&user.LastName,
	)

	if user.ID.Valid {
		shop.CreatedBy = &common.ShortUser{
			Id:        user.ID.String,
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "error while getting shop")
	}

	return &shop, nil
}

func (p *shopRepo) Update(entity *corporate_service.UpdateShopRequest) (string, error) {

	var query = `
		UPDATE	
			"shop"
		SET
			title = $2,
			phone_number = $3,
			size = $4,
			address = $5,
			description = $6
		WHERE
			id = $1
		RETURNING 
			id;
	`
	err := p.db.QueryRow(query, entity.Id, entity.Title, entity.PhoneNumber, entity.Size, entity.Address, entity.Description).Scan(&entity.Id)
	if err != nil {
		return "", errors.Wrap(err, "error while updating shop")
	}

	return entity.Id, nil
}

func (p *shopRepo) GetAll(req *common.SearchRequest) (*corporate_service.GetAllShopsResponse, error) {
	var (
		res corporate_service.GetAllShopsResponse = corporate_service.GetAllShopsResponse{
			Data: make([]*corporate_service.Shop, 0),
		}
		searchFields map[string]interface{} = map[string]interface{}{
			"limit":      req.Limit,
			"offset":     req.Limit * (req.Page - 1),
			"search":     req.Search,
			"company_id": req.Request.CompanyId,
		}
	)

	filter := `
		WHERE
			sh.deleted_at = 0  and sh.company_id = :company_id
	`

	rQuery := `
		SELECT 
			sh.id,
			sh.title,
			sh.phone_number,
			sh.size,
			sh.address,
			sh.description,
			sh.created_at,
			sh.number_of_cashboxes,
			sh.created_by,
			u.first_name,
			u.last_name
		FROM 
			"shop" AS sh
		LEFT JOIN "user" u ON u.id = sh.created_by AND u.deleted_at = 0
	`

	if req.GetSearch() != "" {
		filter += `
		AND
		(
			sh.title ILIKE '%' || :search || '%'
			OR
			sh.phone_number ILIKE '%' || :search || '%'
		)
	`
	}

	rQuery += filter + `ORDER BY sh.created_at DESC LIMIT :limit OFFSET :offset`

	rows, err := p.db.NamedQuery(rQuery, searchFields)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting shop pagination")
	}

	defer rows.Close()

	for rows.Next() {

		var (
			shop corporate_service.Shop
			user models.NullShortUser
		)

		err = rows.Scan(
			&shop.Id,
			&shop.Title,
			&shop.PhoneNumber,
			&shop.Size,
			&shop.Address,
			&shop.Description,
			&shop.CreatedAt,
			&shop.NumberOfCashbox,
			&user.ID,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error while getting shop pagination rows.Scan")
		}

		if user.ID.Valid {
			shop.CreatedBy = &common.ShortUser{
				Id:        user.ID.String,
				FirstName: user.FirstName.String,
				LastName:  user.LastName.String,
			}
		}

		res.Data = append(res.Data, &shop)
	}

	countQuery := `
		SELECT 
			COUNT(*) AS total
		FROM 
			"shop" AS sh
		LEFT JOIN "user" u ON u.id = sh.created_by AND u.deleted_at = 0
	`
	countQuery += filter

	smtm, err := p.db.PrepareNamed(countQuery)
	if err != nil {
		return nil, errors.Wrap(err, "error while scanning count")
	}

	defer smtm.Close()

	err = smtm.QueryRow(searchFields).Scan(&res.Total)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (p *shopRepo) Delete(entity *common.RequestID) (string, error) {

	var query = `
		UPDATE
			"shop" 
		SET
			deleted_at=extract(epoch from now())::bigint 
		WHERE
			id = $1 AND deleted_at = 0
	`

	res, err := p.db.Exec(query, entity.Id)
	if err != nil {
		return "", errors.Wrap(err, "error while deleting shop")
	}

	i, err := res.RowsAffected()
	if err != nil {
		return "", errors.Wrap(err, "error while deleting shop")
	}

	if i == 0 {
		return "", errors.New("shop not found")
	}

	return entity.Id, nil
}
