package main

import (
	"fmt"

	"dasolutions.sk/go-srui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload()

	page := func(ctx *ui.Context) string {
		return app.Html("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
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
