package main

import (
	"Simple-Golang-Web/internal/handler"
	"Simple-Golang-Web/views/layouts"
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/static/{path...}", func(c *core.RequestEvent) error {
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(c.Response, c.Request)
			return nil
		})

		e.Router.GET("/", func(c *core.RequestEvent) error {
			data := []string{
				"Anderson Tirza Liman",
				"Angelo Benhanan Abinaya Fuun",
				"Jovanus Irwan Susanto",
				"Raymundo Rafaelito Maryos Von Woloblo",
			}

			return handler.Render(c, http.StatusOK, layouts.BaseLayout(data))
		})
		return e.Next()
	})

	err := app.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
