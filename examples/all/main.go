package main

import (
	"github.com/michalCapo/go-srui/ui"
)

// simple registry of routes for menu rendering
type route struct {
	Path  string
	Title string
}

var routes = []route{
	{Path: "/", Title: "Home"},
	{Path: "/hello", Title: "Hello"},
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
					ui.Div("font-bold mr-4")("go-srui: examples"),
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

	// Home page lists links
	app.Page("/", layout("Examples Home", func(ctx *ui.Context) string {
		return ui.Div("grid gap-2")(
			ui.Div("text-xl font-bold")("Choose an example:"),
			ui.Map(routes[1:], func(r *route, _ int) string {
				return ui.Div("")(
					ui.A("text-blue-700 underline", ui.Href(r.Path), ctx.Load(r.Path))(r.Title),
				)
			}),
		)
	}))

	// Individual example pages
	app.Page("/hello", layout("Hello", HelloContent))
	app.Page("/button", layout("Button", ButtonContent))
	app.Page("/counter", layout("Counter", CounterContent))
	app.Page("/login", layout("Login", LoginContent))
	app.Page("/showcase", layout("Showcase", ShowcaseContent))

	app.Listen(":1422")
}
