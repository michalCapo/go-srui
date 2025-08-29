package main

import (
	"fmt"

	"github.com/michalCapo/go-srui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload(true)

	page := func(ctx *ui.Context) string {
		return app.HTML("Test", "p-8 bg-gray-200",
			ui.Div("flex flex-row gap-4")(
				// another use of components rendefing using struct with method
				Counter(3).Render(ctx),
			),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}

// we want to create conter with initial count value
func Counter(count int) *TCounter {
	return &TCounter{Count: count}
}

// struct definition
type TCounter struct {
	Count int
}

func (counter *TCounter) Increment(ctx *ui.Context) string {
	// scan request body and update counter struct
	ctx.Body(counter)

	// inscrese count by 1
	counter.Count++

	// render component
	return counter.Render(ctx)
}

func (counter *TCounter) Decrement(ctx *ui.Context) string {
	// scan body
	ctx.Body(counter)

	counter.Count--

	if counter.Count < 0 {
		counter.Count = 0
	}

	// render component as result of this action
	return counter.Render(ctx)
}

func (counter *TCounter) Render(ctx *ui.Context) string {
	// temporary id
	target := ui.Target()

	// renger html, see target (plachodler) at the end, this is place where action result will be rendered
	return ui.Div("flex gap-2 items-center bg-purple-500 rounded text-white p-px", target)(
		ui.Button().
			// click will call decrement action with counter variable as values sent to this action and result will be rendered at target place
			Click(ctx.Call(counter.Decrement, counter).Replace(target)).
			Class("rounded-l px-5").
			Render("-"),

		// display current count
		ui.Div("text-2xl")(fmt.Sprintf("%d", counter.Count)),

		ui.Button().
			// with action result you can replace (overwrite the target component) or render (inline into target component)
			Click(ctx.Call(counter.Increment, counter).Replace(target)).
			Class("rounded-r px-5").
			Render("+"),
	)
}
