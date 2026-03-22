package handler

import (
	"github.com/a-h/templ"
	"github.com/pocketbase/pocketbase/core"
)

func Render(c *core.RequestEvent, status int, t templ.Component) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	return t.Render(c.Request.Context(), c.Response)
}
