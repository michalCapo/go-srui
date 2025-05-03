package main

import (
	"dasolutions.sk/goui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload()

	buttonId := ui.Target()
	var show **ui.Callable

	button := func(ctx *ui.Context) string {
		return ui.Button(buttonId).
			Click(ctx.Call(show).Replace(buttonId)).
			Class("rounded").
			Color(ui.Blue).
			Render("Click me")
	}

	hide := app.Callable(button)

	page := func(ctx *ui.Context) string {
		show = ctx.Callable(func(ctx *ui.Context) string {
			return ui.Div("flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4", buttonId)(
				"Clicked",
				ui.Button().
					Click(ctx.Call(hide).Replace(buttonId)).
					Class("rounded").
					Color(ui.Red).
					Render("Hide me"),
			)
		})

		return app.Html("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					button(ctx),
				),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}
