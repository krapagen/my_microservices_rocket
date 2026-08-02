package user

import (
	"context"

	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
)

func (r *repository) GetByUUID(ctx context.Context, inputUuid string) (model.User, error) {
	return r.getBy(
		ctx,
		"iam/internal/repository/user/GetByUUID",
		"SELECT uuid, login, password_hash, created_at, updated_at FROM users WHERE uuid = $1",
		"uuid",
		inputUuid,
	)
}
