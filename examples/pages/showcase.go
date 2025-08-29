package pages

import (
    "github.com/go-playground/validator/v10"
    "github.com/michalCapo/go-srui/ui"
    "time"
)

// Demo form showcasing inputs and validation
type DemoForm struct {
    Name      string    `validate:"required"`
    Email     string    `validate:"required,email"`
    Phone     string
    Password  string    `validate:"required,min=6"`
    Age       int       `validate:"gte=0,lte=120"`
    Price     float64   `validate:"gte=0"`
    Bio       string
    Gender    string    `validate:"oneof=male female other"`
    Country   string
    Agree     bool      `validate:"eq=true"`
    BirthDate time.Time
    AlarmTime time.Time
    Meeting   time.Time
}

var (
    demoTarget = ui.Target()
)

func (f *DemoForm) Submit(ctx *ui.Context) string {
    if err := ctx.Body(f); err != nil {
        return f.Render(ctx, &err)
    }
    v := validator.New()
    if err := v.Struct(f); err != nil {
        return f.Render(ctx, &err)
    }
    ctx.Success("Form submitted successfully")
    return f.Render(ctx, nil)
}

func (f *DemoForm) Render(ctx *ui.Context, err *error) string {
    countries := ui.MakeOptions([]string{"", "USA", "Slovakia", "Germany", "Japan"})
    genders := []ui.AOption{{ID: "male", Value: "Male"}, {ID: "female", Value: "Female"}, {ID: "other", Value: "Other"}}

    return ui.Div("grid gap-4 sm:gap-6 lg:grid-cols-2 items-start w-full", demoTarget)(
        ui.Form("flex flex-col gap-4 bg-white p-6 rounded-lg shadow w-full", demoTarget, ctx.Submit(f.Submit).Replace(demoTarget))(
            ui.Div("text-xl font-bold")("Component Showcase Form"),
            ui.ErrorForm(err, nil),

            ui.IText("Name", f).Required().Render("Name"),
            ui.IEmail("Email", f).Required().Render("Email"),
            ui.IPhone("Phone", f).Render("Phone"),
            ui.IPassword("Password").Required().Render("Password"),

            ui.INumber("Age", f).Numbers(0, 120, 1).Render("Age"),
            ui.INumber("Price", f).Format("%.2f").Render("Price (USD)"),
            ui.IArea("Bio", f).Rows(4).Render("Short Bio"),

            ui.Div("block sm:hidden")(
                ui.Div("text-sm font-bold")("Gender"),
                ui.IRadio("Gender", f).Value("male").Render("Male"),
                ui.IRadio("Gender", f).Value("female").Render("Female"),
                ui.IRadio("Gender", f).Value("other").Render("Other"),
            ),
            ui.Div("hidden sm:block overflow-x-auto")(
                ui.IRadioButtons("Gender", f).Options(genders).Render("Gender"),
            ),
            ui.ISelect("Country", f).Options(countries).Render("Country"),
            ui.ICheckbox("Agree", f).Required().Render("I agree to the terms"),

            ui.IDate("BirthDate", f).Render("Birth Date"),
            ui.ITime("AlarmTime", f).Render("Alarm Time"),
            ui.IDateTime("Meeting", f).Render("Meeting (Local)"),

            ui.Div("flex gap-2 mt-2")(
                ui.Button().Submit().Color(ui.Blue).Class("rounded").Render("Submit"),
                ui.Button().Reset().Color(ui.Gray).Class("rounded").Render("Reset"),
            ),
        ),

        ui.Div("flex flex-col gap-4 w-full")(
            ui.Div("bg-white p-6 rounded-lg shadow flex flex-col gap-2 w-full")(
                ui.Div("text-xl font-bold")("Buttons & Colors"),
                ui.Div("grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2")(
                    ui.Button().Color(ui.Blue).Class("rounded w-full").Render("Blue"),
                    ui.Button().Color(ui.Green).Class("rounded w-full").Render("Green"),
                    ui.Button().Color(ui.Red).Class("rounded w-full").Render("Red"),
                    ui.Button().Color(ui.Purple).Class("rounded w-full").Render("Purple"),
                    ui.Button().Color(ui.Yellow).Class("rounded w-full").Render("Yellow"),
                    ui.Button().Color(ui.Gray).Class("rounded w-full").Render("Gray"),
                ),
            ),

            ui.Div("bg-white p-6 rounded-lg shadow flex flex-col gap-3 w-full")(
                ui.Div("text-xl font-bold")("Counter (Actions)"),
                (&TCounter{Count: 2}).Render(ctx),
            ),

            ui.Div("bg-white p-6 rounded-lg shadow flex flex col gap-3 w-full")(
                ui.Div("text-xl font-bold")("Simple Table"),
                func() string {
                    cards := ui.Div("space-y-2 sm:hidden")(
                        ui.Div("border rounded p-3 flex justify-between")(
                            ui.Div("")(
                                ui.Div("text-sm text-gray-500")("Name"),
                                ui.Div("font-semibold")("Alice"),
                            ),
                            ui.Div("text-right")(
                                ui.Div("text-sm text-gray-500")("Country"),
                                ui.Div("")("USA"),
                                ui.Div("text-sm text-gray-500 mt-1")("Age: 29"),
                            ),
                        ),
                        ui.Div("border rounded p-3 flex justify-between")(
                            ui.Div("")(
                                ui.Div("text-sm text-gray-500")("Name"),
                                ui.Div("font-semibold")("Bob"),
                            ),
                            ui.Div("text-right")(
                                ui.Div("text-sm text-gray-500")("Country"),
                                ui.Div("")("Germany"),
                                ui.Div("text-sm text-gray-500 mt-1")("Age: 35"),
                            ),
                        ),
                        ui.Div("border rounded p-3 flex justify-between")(
                            ui.Div("")(
                                ui.Div("text-sm text-gray-500")("Name"),
                                ui.Div("font-semibold")("Miro"),
                            ),
                            ui.Div("text-right")(
                                ui.Div("text-sm text-gray-500")("Country"),
                                ui.Div("")("Slovakia"),
                                ui.Div("text-sm text-gray-500 mt-1")("Age: 41"),
                            ),
                        ),
                    )

                    t := ui.SimpleTable(3, "hidden sm:table w-full text-left border-collapse table-fixed text-sm whitespace-normal break-words")
                    t.Class(0, "font-bold")
                    t.Field("Name").Field("Country").Field("Age")
                    t.Field("Alice").Field("USA").Field("29")
                    t.Field("Bob").Field("Germany").Field("35")
                    t.Field("Miro").Field("Slovakia").Field("41")

                    return ui.Div("")(cards + ui.Div("overflow-x-auto sm:overflow-visible")(t.Render()))
                }(),
            ),

            ui.Div("bg-white p-6 rounded-lg shadow flex flex-col gap-3 w-full")(
                ui.Div("text-xl font-bold")("Markdown"),
                ui.Markdown("prose prose-sm sm:prose max-w-none")(`# Heading\n\n- Item 1\n- Item 2\n\n**Bold** and _italic_.`),
            ),

            ui.Div("bg-white p-6 rounded-lg shadow flex flex-col gap-3 w-full")(
                ui.Div("text-xl font-bold")("Client CAPTCHA (demo)"),
                ui.Div("w-full overflow-x-auto")(ui.Captcha2()),
            ),
        ),
    )
}

func ShowcaseContent(ctx *ui.Context) string {
    form := &DemoForm{}
    return ui.Div("max-w-full sm:max-w-6xl mx-auto flex flex-col gap-6 w-full")(
        ui.Div("text-3xl font-bold")("Go-SRUI Component Showcase"),
        form.Render(ctx, nil),
    )
}

