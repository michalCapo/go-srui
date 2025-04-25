package main

import (
	"fmt"

	"dasolutions.sk/goui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload()

	buttonId := ui.Target()
	var show **ui.Method

	button := func(ctx *ui.Context) string {
		return ui.Button(buttonId).
			Click(ctx.Call(show).Replace(buttonId)).
			Class("rounded").
			Color(ui.Blue).
			Render("Click me")
	}

	hide := app.Callable(button)

	page := func(ctx *ui.Context) string {
		open := func(ctx *ui.Context) string {
			return ui.Div("flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4", buttonId)(
				"Clicked",
				ui.Button().
					Click(ctx.Call(hide).Replace(buttonId)).
					Class("rounded").
					Color(ui.Red).
					Render("Hide me"),
			)
		}

		show = ctx.Callable(open)

		return app.Html("Test", "p-8",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					"Hello",
					button(ctx),
				),

				Counter().Render(ctx, 0),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}

func Counter() *TCounter {
	return &TCounter{}
}

type TCounter struct {
	Count int
}

func (counter *TCounter) Increment(ctx *ui.Context) string {
	data := &TCounter{}
	ctx.Body(data)

	return counter.Render(ctx, data.Count+1)
}

func (counter *TCounter) Decrement(ctx *ui.Context) string {
	data := &TCounter{}
	ctx.Body(data)

	return counter.Render(ctx, data.Count-1)
}

func (counter *TCounter) Render(ctx *ui.Context, count int) string {
	target := ui.Target()
	up := ctx.Callable(counter.Increment)
	down := ctx.Callable(counter.Decrement)

	return ui.Div("flex gap-2 items-center bg-purple-500 rounded text-white p-px pl-4", target)(
		ui.Div("text-2xl")(fmt.Sprintf("%d", count)),

		ui.Button().
			Click(ctx.Call(up, TCounter{Count: count}).Replace(target)).
			Class("rounded").
			Render("Increment"),

		ui.Button().
			Click(ctx.Call(down, TCounter{Count: count}).Replace(target)).
			Class("rounded").
			Render("Decrement"),
	)
}
