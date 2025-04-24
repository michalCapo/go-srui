package main

import (
	"dasolutions.sk/test/ui"
)

func main() {
	app := ui.MakeApp("sk")
	app.Autoreload()

	buttonId := ui.Target()
	var show **ui.Method

	onHide := func(ctx *ui.Context) string {
		return ui.Button(buttonId).
			Click(ctx.Call(show).Replace(buttonId)).
			Class("rounded").
			Color(ui.Blue).
			Render("Click me")
	}

	page := func(ctx *ui.Context) string {
		hide := ctx.Callable(onHide)

		onClick := func(ctx *ui.Context) string {
			return ui.Div("flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4", buttonId)(
				"Clicked",
				ui.Button().
					Click(ctx.Call(hide).Replace(buttonId)).
					Class("rounded").
					Color(ui.Red).
					Render("Hide me"),
			)
		}

		show = ctx.Callable(onClick)

		return app.Html("Test", "p-8",
			ui.Div("flex justify-start gap-4 items-center")(
				ui.Div("")("Hello"),
				ui.Button(buttonId).
					Click(ctx.Call(show).Replace(buttonId)).
					Class("rounded").
					Color(ui.Blue).
					Render("Click me"),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}
