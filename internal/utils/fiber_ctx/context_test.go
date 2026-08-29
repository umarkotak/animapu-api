package fiber_ctx

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestWrapPreservesResponseHeadersAndBodies(t *testing.T) {
	app := fiber.New()
	app.Get("/json", Wrap(func(c *Context) {
		c.Header("X-Test", "json")
		c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
	}))
	app.Get("/raw", Wrap(func(c *Context) {
		c.Writer.Header().Set("X-Test", "raw")
		c.Writer.WriteHeader(http.StatusAccepted)
		_, _ = c.Writer.Write([]byte("ok"))
	}))

	for _, expected := range []struct {
		path   string
		status int
		header string
		body   string
	}{
		{"/json", http.StatusCreated, "json", `{"status":"ok"}`},
		{"/raw", http.StatusAccepted, "raw", "ok"},
	} {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, expected.path, nil))
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != expected.status || resp.Header.Get("X-Test") != expected.header || strings.TrimSpace(string(body)) != expected.body {
			t.Fatalf("%s: status=%d header=%q body=%q", expected.path, resp.StatusCode, resp.Header.Get("X-Test"), body)
		}
	}
}
