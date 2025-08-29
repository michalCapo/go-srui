package main

import (
	"github.com/michalCapo/go-srui/examples/pages"
	"github.com/michalCapo/go-srui/ui"
)

// simple registry of routes for menu rendering
type route struct {
	Path  string
	Title string
}

var routes = []route{
	{Path: "/", Title: "Hello"},
	{Path: "/button", Title: "Button"},
	{Path: "/counter", Title: "Counter"},
	{Path: "/login", Title: "Login"},
	{Path: "/showcase", Title: "Showcase"},
}

func main() {
	app := ui.MakeApp("en")
	app.Autoreload(true)

	// layout builder with top menu
	layout := func(title string, body func(*ui.Context) string) ui.Callable {
		return func(ctx *ui.Context) string {
			// simple navigation
			nav := ui.Div("bg-white shadow mb-6")(
				ui.Div("max-w-5xl mx-auto px-4 py-3 flex flex-wrap gap-2 items-center")(
					ui.Div("flex flex-wrap gap-2")(
						ui.Map(routes, func(r *route, _ int) string {
							return ui.A("px-3 py-1 rounded hover:bg-gray-200",
								ui.Href(r.Path),
								ctx.Load(r.Path),
							)(r.Title)
						}),
					),
				),
			)

			content := body(ctx)
			return app.HTML(title, "p-4 bg-gray-200 min-h-screen", nav+ui.Div("max-w-5xl mx-auto px-2")(content))
		}
	}

	// Individual example pages
	app.Page("/", layout("Hello", pages.HelloContent))
	app.Page("/button", layout("Button", pages.ButtonContent))
	app.Page("/counter", layout("Counter", pages.CounterContent))
	app.Page("/login", layout("Login", pages.LoginContent))
	app.Page("/showcase", layout("Showcase", pages.ShowcaseContent))

	app.Listen(":1422")
}
