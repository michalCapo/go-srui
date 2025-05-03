package main

import (
	"dasolutions.sk/goui/ui"
	"github.com/go-playground/validator/v10"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload()

	page := func(ctx *ui.Context) string {
		return app.Html("Test", "p-8 bg-gray-200",
			LoginForm("user").Render(ctx, nil),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}

func LoginForm(name string) *TLoginForm {
	return &TLoginForm{Name: name}
}

type TLoginForm struct {
	Name     string `validate:"required,oneof=user"`
	Password string `validate:"required,oneof=password"`
}

func (form *TLoginForm) Success(ctx *ui.Context) string {
	return ui.Div("text-green-600 max-w-md p-8 text-center font-bold rounded-lg bg-white shadow-xl")("Success")
}

func (form *TLoginForm) Login(ctx *ui.Context) string {
	if err := ctx.Body(form); err != nil {
		return form.Render(ctx, &err)
	}

	v := validator.New()
	if err := v.Struct(form); err != nil {
		return form.Render(ctx, &err)
	}

	return form.Success(ctx)
}

func (form *TLoginForm) Render(ctx *ui.Context, err *error) string {
	var Translations = map[string]string{
		"Name":              "User name",
		"has invalid value": "is invalid",
	}

	target := ui.Target()
	login := ctx.Callable(form.Login)

	return ui.Form("flex flex-col gap-4 max-w-md bg-white p-8 rounded-lg shadow-xl", target, ctx.Submit(login).Replace(target))(
		ui.ErrorForm(err, &Translations),

		ui.IText("Name", form).
			Required().
			Error(err).
			Render("Name"),

		ui.IPassword("Password").
			Required().
			Error(err).
			Render("Password"),

		ui.Button().
			Submit().
			Color(ui.Blue).
			Class("rounded").
			Render("Login"),
	)
}
