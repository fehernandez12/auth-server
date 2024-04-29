package mapper

import "auth-server/models"

func UserToUserDto(user *models.User) *models.UserDto {
	return &models.UserDto{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Enabled:  user.Enabled,
		Created:  user.CreatedAt.String(),
		Updated:  user.UpdatedAt.String(),
		Roles:    RolesToRoleDtos(user.Roles),
	}
}

func UsersToUserDtos(users []*models.User) []*models.UserDto {
	userDtos := make([]*models.UserDto, 0)
	for _, user := range users {
		userDtos = append(userDtos, UserToUserDto(user))
	}
	return userDtos
}
