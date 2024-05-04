package mapper

import (
	"auth-server/models"
)

func ClientsToClientDtos(clients []*models.Client) []*models.ClientDto {
	var dtos []*models.ClientDto
	for _, client := range clients {
		dtos = append(dtos, ClientToClientDto(client))
	}
	return dtos
}

func ClientToClientDto(client *models.Client) *models.ClientDto {
	return &models.ClientDto{
		ID:          client.ID.String(),
		ClientName:  client.ClientName,
		RedirectURI: client.RedirectURI,
	}
}
