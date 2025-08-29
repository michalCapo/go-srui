package pages

import "github.com/michalCapo/go-srui/ui"

func ButtonContent(ctx *ui.Context) string {
    buttonId := ui.Target()
    var show ui.Callable

    button := func(ctx *ui.Context) string {
        return ui.Button(buttonId).
            Click(ctx.Call(show).Replace(buttonId)).
            Class("rounded").
            Color(ui.Blue).
            Render("Click me")
    }

    show = func(ctx *ui.Context) string {
        return ui.Div("flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4", buttonId)(
            "Clicked",
            ui.Button().
                Click(ctx.Call(button).Replace(buttonId)).
                Class("rounded").
                Color(ui.Red).
                Render("Hide me"),
        )
    }

    return ui.Div("flex flex-row gap-4")(
        ui.Div("flex justify-start gap-4 items-center")(button(ctx)),
    )
}

