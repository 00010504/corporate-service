package postgres

import (
	"genproto/common"

	"github.com/Invan2/invan_corporate_service/pkg/logger"
	"github.com/Invan2/invan_corporate_service/storage/repo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type userRepo struct {
	log logger.Logger
	db  *sqlx.DB
}

func NewUserRepo(log logger.Logger, db *sqlx.DB) repo.UserI {
	return &userRepo{
		log: log,
		db:  db,
	}
}

func (u *userRepo) Upsert(entity *common.UserCreatedModel) error {

	query := `
		INSERT INTO "user" (
			"id",
			"first_name",
			"last_name",
			"phone_number",
			"user_type_id"
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) ON CONFLICT ("id") 
		DO 
			UPDATE SET
				"first_name" = $2,
				"last_name" = $3,
				"phone_number" = $4,
				"user_type_id" = $5
	`

	if _, err := u.db.Exec(query, entity.Id, entity.FirstName, entity.LastName, entity.PhoneNumber, entity.UserTypeId); err != nil {
		return errors.Wrap(err, "error while upsert user")
	}

	return nil
}

func (u *userRepo) Delete(id string) error {
	query := `
		UPDATE "user" SET "deleted_at"=extract(epoch from now())::bigint WHERE "id" =  $1
	`

	_, err := u.db.Exec(query, id)
	if err != nil {
		return errors.Wrap(err, "error while delete user")
	}

	return nil
}
