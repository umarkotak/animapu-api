package anime_controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
)

func TestUpdateAnimeTagsRejectsNonAdmin(t *testing.T) {
	app := fiber.New()
	app.Put("/admin/animes/:anime_source/:source_id/tags", fiber_ctx.Wrap(UpdateAnimeTagsBySource))
	request := httptest.NewRequest(http.MethodPut, "/admin/animes/test/title/tags", strings.NewReader(`{"mode":"add","tags":["for_kids:true"]}`))
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != models.ERROR_MAP[models.ErrUnauthorized].StatusCode {
		t.Fatalf("status = %d; want unauthorized", response.StatusCode)
	}
}
