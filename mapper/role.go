package mapper

import "auth-server/models"

func RoleDtoToRole(roleDto *models.RoleDto) *models.Role {
	permissions := make([]*models.Permission, 0)
	for _, permission := range roleDto.Permissions {
		permissions = append(permissions, &models.Permission{
			Name: permission,
		})
	}
	return &models.Role{
		Name:        roleDto.Name,
		Permissions: permissions,
	}
}

func RoleToRoleDto(role *models.Role) *models.RoleDto {
	permissions := make([]string, 0)
	for _, permission := range role.Permissions {
		permissions = append(permissions, permission.Name)
	}
	return &models.RoleDto{
		Name:        role.Name,
		Permissions: permissions,
	}
}
