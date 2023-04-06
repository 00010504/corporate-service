package postgres

import (
	"encoding/json"
	"genproto/common"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type companySizeRepo struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewCompanySizeRepo(log logger.Logger, db *sqlx.DB) repo.CompanySizeI {
	return &companySizeRepo{
		db:  db,
		log: log,
	}
}

func (c *companySizeRepo) Create(entity *corporate_service.CreateCompanySizeRequest) (*common.ResponseID, error) {

	return nil, nil

}

func (c *companySizeRepo) GetAll(entity *common.SearchRequest) (*corporate_service.GetAllCompanySizeResponse, error) {

	res := corporate_service.GetAllCompanySizeResponse{
		Data: make([]*corporate_service.ShortCompanySize, 0),
	}

	values := map[string]interface{}{
		"limit":  entity.Limit,
		"offset": (entity.Page - 1) * entity.Limit,
		"search": entity.Search,
	}

	query := `
		SELECT
			"id",
			"name_tr",
			"description",
			"from",
			"to"
		FROM "company_size"
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
			companySize corporate_service.ShortCompanySize

			name []byte

			description []byte
		)

		if err := rows.Scan(&companySize.Id, &name, &description, &companySize.From, &companySize.To); err != nil {
			return nil, errors.Wrap(err, "error while scan company size")
		}

		if len(name) > 0 {
			if err := json.Unmarshal(name, &companySize.Name); err != nil {
				return nil, errors.Wrap(err, "error while unmarshaling company size name")
			}
		}

		if len(description) > 0 {
			if err := json.Unmarshal(description, &companySize.Description); err != nil {
				return nil, errors.Wrap(err, "error while unmarshaling company size description")
			}
		}

		res.Data = append(res.Data, &companySize)
	}

	queryCount := `
		SELECT
			count(1)
		FROM "company_size"
		` + filter

	stmt, err := c.db.PrepareNamed(queryCount)
	if err != nil {
		return nil, err
	}

	if err := stmt.QueryRow(values).Scan(&res.Total); err != nil {
		return nil, errors.Wrap(err, "error while scan total count of company sizes")
	}

	return &res, nil
}
