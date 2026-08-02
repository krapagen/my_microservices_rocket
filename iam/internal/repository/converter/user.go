package converter

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	"github.com/krapagen/my_microservices_rocket/iam/internal/repository/record"
)

func ToUserRecord(user model.User) record.UserRecord {
	recordUser := record.UserRecord{
		UUID:         user.UUID.String(),
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
	}
	if user.UpdatedAt != nil {
		recordUser.UpdatedAt = user.UpdatedAt
	}
	return recordUser
}

func ToModelUser(userRecord record.UserRecord) (model.User, error) {
	uuidModel, err := uuid.Parse(userRecord.UUID)
	if err != nil {
		return model.User{}, fmt.Errorf("ошибка парсинга UUID пользователя: %w", errs.ErrInvalidUUID)
	}
	modelUser := model.User{
		UUID:         uuidModel,
		Login:        userRecord.Login,
		PasswordHash: userRecord.PasswordHash,
		CreatedAt:    userRecord.CreatedAt,
	}
	if userRecord.UpdatedAt != nil {
		modelUser.UpdatedAt = userRecord.UpdatedAt
	}
	return modelUser, nil
}
