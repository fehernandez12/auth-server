package mapper

import "auth-server/models"

func RoleDtoToRole(roleDto *models.RoleDto) *models.Role {
	return &models.Role{
		Name: roleDto.Name,
	}
}

func RoleToRoleDto(role *models.Role) *models.RoleDto {
	return &models.RoleDto{
		ID:          role.ID.String(),
		Name:        role.Name,
		Application: ApplicationToApplicationDto(role.Application),
	}
}

func RolesToRoleDtos(roles []*models.Role) []*models.RoleDto {
	roleDtos := make([]*models.RoleDto, 0)
	for _, role := range roles {
		roleDtos = append(roleDtos, RoleToRoleDto(role))
	}
	return roleDtos
}
