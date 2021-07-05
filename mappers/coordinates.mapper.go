package mappers

import (
	"github.com/jental/webbugs-common-go/contracts"
	"github.com/jental/webbugs-common-go/models"
)

func MapFromCoordinates(crd *models.Coordinates) contracts.CoordinatesContract {
	return contracts.CoordinatesContract{
		X: crd.X,
		Y: crd.Y,
		Z: crd.Z,
	}
}

func MapToCoordinates(crd *contracts.CoordinatesContract) models.Coordinates {
	return models.NewCoordinates(crd.X, crd.Y, crd.Z)
}
