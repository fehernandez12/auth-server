package mapper

import "auth-server/models"

func PermissionsToPermissionDtos(permissions []*models.Permission) []*models.PermissionDto {
	permissionDtos := make([]*models.PermissionDto, 0)
	for _, permission := range permissions {
		permissionDtos = append(permissionDtos, PermissionToPermissionDto(permission))
	}
	return permissionDtos
}

func PermissionToPermissionDto(permission *models.Permission) *models.PermissionDto {
	return &models.PermissionDto{
		ID:   permission.ID,
		Name: permission.Name,
		Role: RoleToRoleDto(permission.Role),
	}
}
