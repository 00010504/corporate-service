package repo

import "genproto/common"

type UserI interface {
	Upsert(entity *common.UserCreatedModel) error
}
