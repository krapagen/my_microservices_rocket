package converter

import (
	"github.com/krapagen/my_microservices_rocket/iam/internal/model"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
)

// UserModelToHTTPDTO конвертирует доменную модель пользователя в OpenAPI HTTP DTO
func UserModelToHTTPDTO(user model.User) *iamv1.UserDto {
	updatedAt := user.CreatedAt
	if user.UpdatedAt != nil {
		updatedAt = *user.UpdatedAt
	}

	return &iamv1.UserDto{
		UUID:      user.UUID,
		Login:     user.Login,
		CreatedAt: user.CreatedAt,
		UpdatedAt: updatedAt,
	}
}
