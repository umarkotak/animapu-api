package anime_controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
)

func TestGetCapabilities(t *testing.T) {
	app := fiber.New()
	app.Get("/animes/:anime_source/capabilities", fiber_ctx.Wrap(GetCapabilities))
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/animes/kuramanime/capabilities", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var body struct {
		Data []contract.AnimeCapability `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || len(body.Data) != 1 || !body.Data[0].GetLatest {
		t.Fatalf("status = %d, capabilities = %#v", response.StatusCode, body.Data)
	}
}

func TestGetLatestRejectsInvalidSpecialFilter(t *testing.T) {
	app := fiber.New()
	app.Get("/animes/:anime_source/latest", fiber_ctx.Wrap(GetLatest))
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/animes/kuramanime/latest?special_filter=invalid", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if want := models.ERROR_MAP[models.ErrInvalidFormat].StatusCode; response.StatusCode != want {
		t.Fatalf("status = %d, want %d", response.StatusCode, want)
	}
}
