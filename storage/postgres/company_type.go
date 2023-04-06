package postgres

import (
	"encoding/json"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type companyTypeRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewCompanyType(log logger.Logger, db *sqlx.DB) repo.CompanyTypeI {
	return &companyTypeRepo{
		db:  db,
		log: log,
	}
}

// /////Create
func (p *companyTypeRepo) Create(entity *corporate_service.CreateCompanyTypeRequest) (*common.ResponseID, error) {

	var query = `
		INSERT INTO
				"company_type" 
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
	id := uuid.New().String()

	_, err := p.db.Exec(
		query,
		id,
		entity.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error while inserting company type")
	}
	return &common.ResponseID{
		Id: id,
	}, nil
}

// ////GetById
func (p *companyTypeRepo) GetById(entity *common.RequestID) (*corporate_service.CompanyType, error) {

	var companyType corporate_service.CompanyType

	var query = `
		SELECT 
			id,
			name, 
			created_at,
			deleted_at
		FROM 
			"company_type"
		WHERE 
			id = $1
	`
	err := p.db.QueryRow(query, entity.Id).Scan(&companyType.Id, &companyType.Name, &companyType.CreatedAt)

	if err != nil {
		return nil, errors.Wrap(err, "error while getting company type")
	}

	return &companyType, nil
}

func (p *companyTypeRepo) Update(entity *corporate_service.UpdateCompanyTypeRequest) (string, error) {
	var query = `
		UPDATE	
			"company_type"
		SET
			name = $2
		WHERE
			id = $1
		RETURNING 
			id;
	`
	err := p.db.QueryRow(query, entity.Id, entity.Name).Scan(&entity.Id)

	if err != nil {
		return "", errors.Wrap(err, "error while updating company type")
	}

	return entity.Id, nil
}

func (p *companyTypeRepo) Delete(entity *common.RequestID) (string, error) {

	var query = `
		UPDATE 
			company_type
		SET 
			deleted_at=extract(epoch from now())::bigint 
		WHERE id = $1 and deleted_at = 0 and company_id = $2`

	res, err := p.db.Exec(query, entity.Id)
	if err != nil {
		return "", errors.Wrap(err, "error while deleting company type")
	}

	i, err := res.RowsAffected()
	if err != nil {
		return "", errors.Wrap(err, "error while deleting company_type")
	}

	if i == 0 {
		return "", errors.New("company_type not found")
	}

	return entity.Id, nil
}

func (c *companyTypeRepo) GetAll(entity *common.SearchRequest) (*corporate_service.GetAllCompanyTypeResponse, error) {

	res := corporate_service.GetAllCompanyTypeResponse{
		Data: make([]*corporate_service.ShortCompanyType, 0),
	}

	values := map[string]interface{}{
		"limit":  entity.Limit,
		"offset": (entity.Page - 1) * entity.Limit,
		"search": entity.Search,
	}

	query := `
			SELECT
				"id",
				"name"
			FROM "company_type"
			WHERE "deleted_at" = 0
		`
	filter := ``

	if entity.Search != "" {
		filter += ` AND (
				"name"->en ILIKE '%' || :search || '%' OR
				"name"->uz ILIKE '%' || :search || '%' OR
				"name"->ru ILIKE '%' || :search || '%' 
			)`
	}

	query += filter

	query += ` LIMIT :limit OFFSET :offset`

	rows, err := c.db.NamedQuery(query, values)
	if err != nil {
		return nil, errors.Wrap(err, "error while select company size")
	}

	defer rows.Close()

	for rows.Next() {

		var (
			companyType corporate_service.ShortCompanyType
			name        []byte
		)

		if err := rows.Scan(&companyType.Id, &name); err != nil {
			return nil, errors.Wrap(err, "error while scan company type")
		}

		if len(name) > 0 {
			if err := json.Unmarshal(name, &companyType.Name); err != nil {
				return nil, errors.Wrap(err, "error while unmarshaling company type name")
			}
		}

		res.Data = append(res.Data, &companyType)
	}

	queryCount := `
			SELECT
				count(1)
			FROM "company_type"
			` + filter

	stmt, err := c.db.PrepareNamed(queryCount)
	if err != nil {
		return nil, errors.Wrap(err, "error while prepare named query")
	}

	if err := stmt.QueryRow(values).Scan(&res.Total); err != nil {
		return nil, errors.Wrap(err, "error while scan total count of company sizes")
	}

	return &res, nil

}
