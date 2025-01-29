package request

import (
	"net/http"

	"github.com/umarkotak/animapu-api/internal/models"
)

func ReqToUser(req *http.Request) (models.UserFirebase, error) {
	if len(req.Header["Animapu-User-Uid"]) < 1 {
		return models.UserFirebase{}, models.ErrNotFound
	}

	return models.UserFirebase{
		Uid: req.Header["Animapu-User-Uid"][0],
	}, nil
}
