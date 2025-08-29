package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/michalCapo/go-srui/ui"
)

func main() {
	app := ui.MakeApp("en")
	app.Autoreload(true)

	page := func(ctx *ui.Context) string {
		return app.HTML("Test", "p-8 bg-gray-200",
			// rendering login form with context and nil as error
			LoginForm("user").Render(ctx, nil),
		)
	}

	app.Page("/", page)
	app.Listen(":1422")
}

func LoginForm(name string) *TLoginForm {
	return &TLoginForm{Name: name}
}

// defining login form with validations for given fields
type TLoginForm struct {
	Name     string `validate:"required,oneof=user"`
	Password string `validate:"required,oneof=password"`
}

// we want to display success message
func (form *TLoginForm) Success(ctx *ui.Context) string {
	return ui.Div("text-green-600 max-w-md p-8 text-center font-bold rounded-lg bg-white shadow-xl")("Success")
}

// Login action
func (form *TLoginForm) Login(ctx *ui.Context) string {
	// scan request body, if there is an error render with using render method of this component
	if err := ctx.Body(form); err != nil {
		return form.Render(ctx, &err)
	}

	v := validator.New()
	// let's validate our input, and display error if any
	if err := v.Struct(form); err != nil {
		return form.Render(ctx, &err)
	}

	// great a successful login
	return form.Success(ctx)
}

// translations for login form
var translations = map[string]string{
	"Name":              "User name",
	"has invalid value": "is invalid",
}

// temporary id
var target = ui.Target()

func (form *TLoginForm) Render(ctx *ui.Context, err *error) string {

	// submiting form will call login action and result will be rendered to target id
	return ui.Form("flex flex-col gap-4 max-w-md bg-white p-8 rounded-lg shadow-xl", target, ctx.Submit(form.Login).Replace(target))(
		// display all error in one place
		ui.ErrorForm(err, &translations),

		// text component
		ui.IText("Name", form).
			// is requered
			Required().
			// if there is specific error for this field display it
			Error(err).
			Render("Name"),

		// password component
		ui.IPassword("Password").
			Required().
			Error(err).
			Render("Password"),

		// submit button, see submit part on form several lines above
		ui.Button().
			Submit().
			Color(ui.Blue).
			Class("rounded").
			Render("Login"),
	)
}
