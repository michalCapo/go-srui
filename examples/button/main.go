package main

import (
	"github.com/michalCapo/go-srui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload(true)

	// create id (placeholder)
	buttonId := ui.Target()
	var show ui.Callable

	// button component, it can be used as component render function or action, please see references to this variable
	button := func(ctx *ui.Context) string {

		// return basic component
		return ui.Button(buttonId).
			// on button click show method is called and result will be placed in html document marked with buttonId target
			Click(ctx.Call(show).Replace(buttonId)).
			Class("rounded").
			Color(ui.Blue).
			Render("Click me")
	}

	page := func(ctx *ui.Context) string {
		show = func(ctx *ui.Context) string {
			// buttonId is used on serveral places, it mark the spot where action should be rendered
			return ui.Div("flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4", buttonId)(
				"Clicked",
				ui.Button().
					// clicking on button will call button action and result will replace the buttonId placeholder
					Click(ctx.Call(button).Replace(buttonId)).
					Class("rounded").
					Color(ui.Red).
					Render("Hide me"),
			)
		}

		// again, at this point we are rending whole page, so use app.Html function
		return app.HTML("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					// button function is used as render function
					button(ctx),
				),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}
