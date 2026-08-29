package fiber_ctx

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

const contextKey = "animapu.fiber_ctx"

type Context struct {
	Ctx     fiber.Ctx
	Request *http.Request
	Writer  responseWriter
}

type responseWriter struct {
	ctx         fiber.Ctx
	header      http.Header
	status      int
	headersSent bool
}

func From(c fiber.Ctx) *Context {
	if ctx, ok := c.Locals(contextKey).(*Context); ok {
		return ctx
	}

	req := new(http.Request)
	_ = fasthttpadaptor.ConvertRequest(c.RequestCtx(), req, true)
	req = req.WithContext(c.Context())
	ctx := &Context{
		Ctx:     c,
		Request: req,
		Writer: responseWriter{
			ctx:    c,
			header: make(http.Header),
			status: http.StatusOK,
		},
	}
	c.Locals(contextKey, ctx)
	return ctx
}

func Wrap(handler func(*Context)) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := From(c)
		handler(ctx)
		ctx.flushHeaders()
		return nil
	}
}

func (c *Context) Param(name string) string { return c.Ctx.Params(name) }
func (c *Context) Query(name string) string { return c.Ctx.Query(name) }

func (c *Context) Set(key string, value any) { c.Ctx.Locals(key, value) }

func (c *Context) Get(key string) (any, bool) {
	value := c.Ctx.Locals(key)
	return value, value != nil
}

func (c *Context) Header(key, value string) { c.Writer.Header().Set(key, value) }

func (c *Context) BindJSON(value any) error { return c.Ctx.Bind().JSON(value) }

func (c *Context) JSON(status int, value any) {
	c.flushHeaders()
	_ = c.Ctx.Status(status).JSON(value)
}

func (c *Context) flushHeaders() {
	if c.Writer.headersSent {
		return
	}
	for key, values := range c.Writer.header {
		for _, value := range values {
			c.Ctx.Append(key, value)
		}
	}
	c.Writer.headersSent = true
}

func (w *responseWriter) Header() http.Header { return w.header }

func (w *responseWriter) WriteHeader(status int) { w.status = status }

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.headersSent {
		for key, values := range w.header {
			for _, value := range values {
				w.ctx.Append(key, value)
			}
		}
		w.headersSent = true
	}
	return len(body), w.ctx.Status(w.status).Send(body)
}
