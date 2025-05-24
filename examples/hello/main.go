package main

import (
	"github.com/michalCapo/go-srui/ui"
)

func main() {
	// add some basic meta, styling and link to cdn tailwind library
	app := ui.MakeApp("en")
	// add autoreload behavior to html head
	app.Autoreload()

	page := func(ctx *ui.Context) string {
		// as this is a page we need to return full HTML, including head and body
		return app.Html("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					// just some html displaying "hello"
					"Hello",
				),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}
