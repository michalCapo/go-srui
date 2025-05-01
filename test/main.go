package main

import (
	"fmt"

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

		return app.Html("Test", "p-8",
			ui.Div("flex flex-row gap-4")(
				ui.Div("flex justify-start gap-4 items-center")(
					"Hello",
					button(ctx),
				),

				Counter(3).Render(ctx),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}

func Counter(count int) *TCounter {
	return &TCounter{Count: count}
}

type TCounter struct {
	Count int
}

func (counter *TCounter) Increment(ctx *ui.Context) string {
	ctx.Body(counter)

	counter.Count++

	return counter.Render(ctx)
}

func (counter *TCounter) Decrement(ctx *ui.Context) string {
	ctx.Body(counter)

	counter.Count--

	if counter.Count < 0 {
		counter.Count = 0
	}

	return counter.Render(ctx)
}

func (counter *TCounter) Render(ctx *ui.Context) string {
	target := ui.Target()
	up := ctx.Callable(counter.Increment)
	down := ctx.Callable(counter.Decrement)

	return ui.Div("flex gap-2 items-center bg-purple-500 rounded text-white p-px", target)(
		ui.Button().
			Click(ctx.Call(down, counter).Replace(target)).
			Class("rounded-l px-5").
			Render("-"),

		ui.Div("text-2xl")(fmt.Sprintf("%d", counter.Count)),

		ui.Button().
			Click(ctx.Call(up, counter).Replace(target)).
			Class("rounded-r px-5").
			Render("+"),
	)
}
