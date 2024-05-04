package mapper

import (
	"auth-server/models"
)

func ApplicationsToApplicationDtos(apps []*models.Application) []*models.ApplicationDto {
	var dtos []*models.ApplicationDto
	for _, app := range apps {
		dtos = append(dtos, &models.ApplicationDto{
			ID:      app.ID.String(),
			AppName: app.AppName,
		})
	}
	return dtos
}

func ApplicationToApplicationDto(app *models.Application) *models.ApplicationDto {
	return &models.ApplicationDto{
		ID:      app.ID.String(),
		AppName: app.AppName,
	}
}
