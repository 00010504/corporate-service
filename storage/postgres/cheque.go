package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"genproto/common"
	"sort"
	"strings"

	"genproto/corporate_service"

	"github.com/Invan2/invan_corporate_service/config"
	"github.com/Invan2/invan_corporate_service/models"
	"github.com/Invan2/invan_corporate_service/pkg/helper"
	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type chequeRepo struct {
	db  *sqlx.DB
	log logger.Logger
	cfg config.Config
}

func NewChequeRepo(db *sqlx.DB, log logger.Logger, cfg config.Config) repo.ChequeI {
	return &chequeRepo{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

func (c *chequeRepo) Create(entity *corporate_service.CreateChequeRequest) (*common.ResponseID, error) {

	tr, err := c.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while staring transaction on creating cheque")
	}

	if strings.Contains(entity.Logo.Image, "http") {
		imageSplittedArr := strings.Split(entity.Logo.Image, "/")
		entity.Logo.Image = imageSplittedArr[len(imageSplittedArr)-1]
	}

	defer func() {

		if err != nil {
			_ = tr.Rollback()
		} else {
			_ = tr.Commit()
		}

	}()

	chequeID := uuid.NewString()

	query := `
		INSERT INTO "cheque" (
			"id",
			"name",
			"message",
			"company_id",
			"created_by"
		) VALUES ($1, $2, $3, $4, $5);
	`

	_, err = tr.Exec(query, chequeID, entity.Name, entity.Message, entity.Request.CompanyId, entity.Request.UserId)
	if err != nil {
		return nil, errors.Wrap(err, "error while insert cheque")
	}

	if entity.Logo == nil {
		entity.Logo = &corporate_service.ChequeLogo{}
	}

	query = `
		INSERT INTO  "cheque_logo" 
			(
				"cheque_id",
				"image", 
				"left", 
				"right",
				"top",	
				"bottom"
			) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tr.Exec(query, chequeID, entity.Logo.Image, entity.Logo.Left, entity.Logo.Right, entity.Logo.Top, entity.Logo.Bottom)
	if err != nil {
		return nil, errors.Wrap(err, "error while insert cheque")
	}

	query = `
		INSERT INTO "cheque_field" (
			"cheque_id",
			"field_id",
			"position",
			"is_added",
			"created_by"
		) VALUES 
	`

	values := make([]interface{}, 0)

	for _, fieldID := range entity.FieldIds {
		query += `(?, ?, ?, ?, ?),`
		values = append(values, chequeID, fieldID.FieldId, fieldID.Position, fieldID.IsAdded, entity.Request.UserId)
	}

	query = strings.TrimSuffix(query, ",")
	query = helper.ReplaceSQL(query, "?")

	query += `
		ON CONFLICT ("cheque_id", "field_id", "position", "deleted_at") DO UPDATE SET position = EXCLUDED.position
	`

	stmt, err := tr.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating cashbox_payment Prepare")
	}

	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return nil, errors.Wrap(err, "error while creating cashbox_payment Exec")
	}

	return &common.ResponseID{
		Id: chequeID,
	}, nil
}

func (c *chequeRepo) Get(entity *common.RequestID) (*corporate_service.Cheque, error) {

	var (
		res corporate_service.Cheque = corporate_service.Cheque{
			Blocks: make([]*corporate_service.RecieptBlock, 0),
		}

		blocksMap map[string]*corporate_service.RecieptBlock = make(map[string]*corporate_service.RecieptBlock)
		user      models.NullShortUser

		logo models.NullChequeLogo
	)

	query := `
		SELECT 
			c."id",
			c."name",
			c.message,
			u.first_name,
			u.last_name,
			u.id,
			l."image",
			l."left",
			l."right",
			l."top",
			l."bottom"
		FROM "cheque" c
		LEFT JOIN "user" u ON u.id = c.created_by AND u.deleted_at = 0
		LEFT JOIN "cheque_logo" l ON l.cheque_id = c.id
		WHERE c.id = $1 AND c.deleted_at = 0 AND c.company_id = $2;
	`

	err := c.db.QueryRow(query, entity.Id, entity.Request.CompanyId).Scan(
		&res.Id,
		&res.Name,
		&res.Message,
		&user.FirstName,
		&user.LastName,
		&user.ID,
		&logo.Image,
		&logo.Left,
		&logo.Right,
		&logo.Top,
		&logo.Bottom,
	)
	if err != nil {
		return nil, err
	}

	if user.ID.Valid {
		res.CreatedBy = &common.ShortUser{
			FirstName: user.FirstName.String,
			LastName:  user.LastName.String,
			Id:        user.ID.String,
		}
	}

	if logo.Image.String != "" {
		res.Logo = &corporate_service.ChequeLogo{
			Image:  fmt.Sprintf("https://%s/%s/%s", c.cfg.MinioEndpoint, config.FileBucketName, logo.Image.String),
			Left:   float32(logo.Left.Float64),
			Right:  float32(logo.Right.Float64),
			Top:    float32(logo.Top.Float64),
			Bottom: float32(logo.Bottom.Float64),
		}
	}

	query = `
		SELECT 
			"id",
			"name",
			"name_tr"
		FROM "receipt_block"
		WHERE deleted_at = 0
		ORDER BY name DESC;
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		block := corporate_service.RecieptBlock{
			Fields: make([]*corporate_service.RecieptField, 0),
		}

		nameTranslation := []byte{}

		if err := rows.Scan(&block.Id, &block.Name, &nameTranslation); err != nil {
			return nil, err
		}

		if len(nameTranslation) > 0 {
			if err := json.Unmarshal(nameTranslation, &block.NameTranslation); err != nil {
				return nil, err
			}
		}

		blocksMap[block.Id] = &block
	}

	rows.Close()

	query = `
		SELECT 
			rf."id",
			rf."name",
			rf."name_tr",
			rf."block_id",
			chf.is_added
		FROM cheque_field chf
		LEFT JOIN receipt_field rf ON rf.id = chf.field_id
		WHERE chf.deleted_at = 0 AND chf.cheque_id = $1
		ORDER BY chf.position ASC;
	`

	rows, err = c.db.Query(query, entity.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting cheque fields. Query")
	}

	defer rows.Close()

	for rows.Next() {
		field := corporate_service.RecieptField{}
		nameTranslation := []byte{}

		var blockID string

		if err := rows.Scan(&field.Id, &field.Name, &nameTranslation, &blockID, &field.IsAdded); err != nil {
			return nil, errors.Wrap(err, "error while getting cheque fields. Scan")
		}

		if len(nameTranslation) > 0 {
			if err := json.Unmarshal(nameTranslation, &field.NameTranslation); err != nil {
				return nil, errors.Wrap(err, "error while getting cheque fields. Unmarshal")
			}
		}

		block, ok := blocksMap[blockID]
		if ok {
			block.Fields = append(block.Fields, &field)
			blocksMap[blockID] = block
		}

	}

	rows.Close()

	for _, block := range blocksMap {
		block.TotalFields = int32(len(block.Fields))
		res.Blocks = append(res.Blocks, block)
	}

	sort.SliceStable(res.Blocks, func(i, j int) bool {
		return res.Blocks[i].Name > res.Blocks[j].Name
	})

	res.TotalBlocks = int32(len(res.Blocks))
	return &res, nil
}

func (c *chequeRepo) GetAll(entity *common.SearchRequest) (*corporate_service.GetAllChequesResponse, error) {

	var (
		res corporate_service.GetAllChequesResponse = corporate_service.GetAllChequesResponse{
			Data: make([]*corporate_service.ShortCheque, 0),
		}

		valuesMap = map[string]interface{}{
			"search":     entity.Search,
			"limit":      entity.Limit,
			"offset":     (entity.Page - 1) * entity.Limit,
			"company_id": entity.Request.CompanyId,
		}
	)

	query := `
		SELECT 
			c."id",
			c."name",
			c."message",
			u.first_name,
			u.last_name,
			u.id
		FROM "cheque" c
		LEFT JOIN "user" u ON u.id = c.created_by AND u.deleted_at = 0
		WHERE  c.deleted_at = 0 AND c.company_id = :company_id
	`

	filter := ``

	if entity.Search != "" {
		filter += `
			AND (
				c.name ILIKE '%' || :search || '%' 
			) 
		`
	}

	query += filter

	query += ` ORDER BY c.created_at DESC LIMIT :limit OFFSET :offset `

	rows, err := c.db.NamedQuery(query, valuesMap)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting cashboxs")
	}

	defer rows.Close()

	for rows.Next() {

		var cheque corporate_service.ShortCheque

		var user models.NullShortUser

		if err := rows.Scan(&cheque.Id, &cheque.Name, &cheque.Message, &user.FirstName, &user.LastName, &user.ID); err != nil {
			return nil, err
		}

		if user.ID.Valid {
			cheque.CreatedBy = &common.ShortUser{
				FirstName: user.FirstName.String,
				LastName:  user.LastName.String,
				Id:        user.ID.String,
			}
		}

		res.Data = append(res.Data, &cheque)

	}

	query = `
		SELECT 
			count(1) as total
		FROM "cheque" c
		LEFT JOIN "user" u ON u.id = c.created_by AND u.deleted_at = 0
		WHERE  c.deleted_at = 0 AND c.company_id = :company_id 
	`

	query += filter

	stmt, err := c.db.PrepareNamed(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while prepare getting total")
	}

	defer stmt.Close()

	if err := stmt.QueryRow(valuesMap).Scan(&res.Total); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *chequeRepo) GetAllRecieptBlock(req *common.SearchRequest) (*corporate_service.GetAllRecieptBlockResponse, error) {
	var (
		res = corporate_service.GetAllRecieptBlockResponse{
			Blocks: make([]*corporate_service.RecieptBlock, 0),
		}

		fieldsMap = make(map[string][]*corporate_service.RecieptField)
	)

	query := `
		SELECT 
			"id",
			"name",
			"name_tr"
		FROM 
			"receipt_block"
		WHERE deleted_at = 0
		ORDER BY name DESC;  
	`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while select blocks")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			block       corporate_service.RecieptBlock
			translation []byte
		)

		err := rows.Scan(&block.Id, &block.Name, &translation)
		if err != nil {
			return nil, err
		}

		if len(translation) > 0 {

			err := json.Unmarshal(translation, &block.NameTranslation)
			if err != nil {
				return nil, err
			}
		}

		fieldsMap[block.Id] = make([]*corporate_service.RecieptField, 0)

		res.Blocks = append(res.Blocks, &block)
	}

	_ = rows.Close()

	query = `
		SELECT
            "id",
            "name",
            "name_tr",
            "block_id"
        FROM receipt_field
        WHERE deleted_at = 0;
	`
	rows, err = c.db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while selecting receipt fields")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			field       corporate_service.RecieptField
			translation []byte
			blockId     string
		)

		err := rows.Scan(&field.Id, &field.Name, &translation, &blockId)
		if err != nil {
			return nil, errors.Wrap(err, "error while scanning receipt field")
		}

		if len(translation) > 0 {
			err := json.Unmarshal(translation, &field.NameTranslation)
			if err != nil {
				return nil, errors.Wrap(err, "error while unmarshaling name translations")

			}
		}

		fieldsMap[blockId] = append(fieldsMap[blockId], &field)
	}

	for _, block := range res.Blocks {
		block.Fields = fieldsMap[block.Id]
	}

	return &res, nil
}

func (r *chequeRepo) Delete(ctx context.Context, req *common.RequestID) (*common.ResponseID, error) {

	query := `
		UPDATE 
			"cheque"
		SET deleted_at=extract(epoch from now())::bigint 
		WHERE "id"=$1 and company_id=$2 and deleted_at=0
	`
	res, err := r.db.Exec(query, req.Id, req.Request.CompanyId)
	if err != nil {
		return nil, err
	}

	if i, _ := res.RowsAffected(); i == 0 {
		return nil, sql.ErrNoRows
	}

	return &common.ResponseID{Id: req.Id}, nil
}

func (r *chequeRepo) Update(ctx context.Context, req *corporate_service.UpdateChequeRequest) (*common.ResponseID, error) {

	tr, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while staring transaction on creating cheque")
	}

	defer func() {

		if err != nil {
			_ = tr.Rollback()
		} else {
			_ = tr.Commit()
		}
	}()

	query := `
	UPDATE  "cheque" SET 
		"name"=$3,
		"message"=$4
		
	WHERE id=$1 AND company_id=$2 AND deleted_at=0
`

	_, err = tr.Exec(query, req.Id, req.Request.CompanyId, req.Name, req.Message)
	if err != nil {
		return nil, errors.Wrap(err, "error while insert cheque")
	}

	if req.Logo == nil {
		req.Logo = &corporate_service.ChequeLogo{}
	}

	if strings.Contains(req.Logo.Image, "http") {
		imageSplittedArr := strings.Split(req.Logo.Image, "/")
		req.Logo.Image = imageSplittedArr[len(imageSplittedArr)-1]
	}

	query = `
		INSERT INTO  "cheque_logo" 
			(
				"cheque_id",
				"image", 
				"left", 
				"right",
				"top",	
				"bottom"
			) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (cheque_id) DO
		UPDATE
			SET	
			"image"=$2, 
			"left"=$3, 
			"right"=$4,
			"top"=$5,	
			"bottom"=$6
			`

	_, err = tr.Exec(query, req.Id, req.Logo.Image, req.Logo.Left, req.Logo.Right, req.Logo.Top, req.Logo.Bottom)
	if err != nil {
		return nil, errors.Wrap(err, "error while insert cheque")
	}

	query = `
		DELETE FROM "cheque_field" 
		WHERE cheque_id = $1 AND deleted_at = 0
	`

	_, err = tr.Exec(query, req.Id)
	if err != nil {
		return nil, errors.Wrap(err, "error while update cheque_field")
	}

	query = `
	INSERT INTO "cheque_field" (
		"cheque_id",
		"field_id",
		"position",
		"is_added",
		"created_by"
	) VALUES 
`

	values := make([]interface{}, 0)

	for _, fieldID := range req.FieldIds {
		query += `(?, ?, ?, ?, ?),`
		values = append(values, req.Id, fieldID.FieldId, fieldID.Position, fieldID.IsAdded, req.Request.UserId)
	}

	query = strings.TrimSuffix(query, ",")
	query = helper.ReplaceSQL(query, "?")

	query += ` ON CONFLICT ("field_id", "cheque_id", "deleted_at", "position") DO UPDATE SET position = EXCLUDED.position, is_added = EXCLUDED.is_added `

	stmt, err := tr.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating cheque_field Prepare")
	}

	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return nil, errors.Wrap(err, "error while updating cheque_field Exec")
	}

	return &common.ResponseID{Id: req.Id}, nil
}
