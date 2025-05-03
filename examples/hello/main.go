package main

import (
	"github.com/michalCapo/go-srui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload()

	page := func(ctx *ui.Context) string {
		return app.Html("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					"Hello",
				),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}
