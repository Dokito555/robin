package http

import (
	"github.com/Dokito555/robin-ums/internal/services"
	"github.com/sirupsen/logrus"
)

type ArtistController struct {
	Log     *logrus.Logger
	Service *services.ArtistService
}

func NewArtistController(log *logrus.Logger, service *services.ArtistService) *ArtistController {
	return &ArtistController{
		Log:     log,
		Service: service,
	}
}