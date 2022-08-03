package request

import (
	"net/http"

	"github.com/umarkotak/animapu-api/internal/models"
)

func ReqToUser(req *http.Request) (models.User, error) {
	if len(req.Header["Animapu-User-Uid"]) < 1 {
		return models.User{}, models.ErrNotFound
	}

	return models.User{
		Uid: req.Header["Animapu-User-Uid"][0],
	}, nil
}
