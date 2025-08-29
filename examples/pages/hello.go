package pages

import "github.com/michalCapo/go-srui/ui"

func HelloContent(ctx *ui.Context) string {
    return ui.Div("flex flex-row gap-4")(
        ui.Div("flex justify-start gap-4 items-center")(
            "Hello",
        ),
    )
}

