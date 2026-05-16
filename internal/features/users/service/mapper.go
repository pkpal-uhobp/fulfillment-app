package users_service

import core_domain "github.com/pkpal-uhobp/fulfillment-app/internal/core/domain"

func toUserDTO(user core_domain.User) UserDTO {
	dto := UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role.String(),
		IsActive:  user.IsActive,
		IsBlocked: user.IsBlocked,
	}
	if user.Phone != nil {
		dto.Phone = *user.Phone
	}
	return dto
}

func toUserDTOs(users []core_domain.User) []UserDTO {
	result := make([]UserDTO, 0, len(users))
	for _, user := range users {
		result = append(result, toUserDTO(user))
	}
	return result
}
