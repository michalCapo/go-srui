![go-srui](https://github.com/user-attachments/assets/ec51e978-8e53-4cb1-8e24-bb5e0494453a)

# Go-SRUI (Server-Rendered UI)

A lightweight, server-rendered UI framework for Go that enables building interactive web applications from your Go code. Define your UI in Go and let the framework handle the rest.

## Features

- No JavaScript required
- No build step required
- No need to learn new syntax for writing HTML
- No plugins required for writing HTML
- Server-side rendering with client-side interactivity
- Built-in form handling and validation
- Automatic WebSocket-based live reload for development
- Tailwind CSS integration
- Session management
- File upload and download support
- Customizable HTML components
- Built-in form validation using go-playground/validator
- Support for various input types (text, email, phone, number, etc.)

## Installation

```bash
go get github.com/michalCapo/go-srui
```

## Quick Start

```go
package main

import "github.com/michalCapo/go-srui/ui"

func main() {
    // Create a new app instance
    app := ui.MakeApp("en")
    
    // Enable live reload for development
    app.Autoreload()

    // Define a page handler
    page := func(ctx *ui.Context) string {
        return app.Html("My App", "p-8 bg-gray-200",
            ui.Div("flex flex-row gap-4")(
                ui.Div("flex justify-start gap-4 items-center")(
                    "Hello, World!",
                ),
            ),
        )
    }

    // Register the page
    app.Page("/", page)
    
    // Start the server
    app.Listen(":8080")
}
```
## How it works

Basically as HTMX, but with Go. You define your components or actions in Go. Register these actions so framework can use them. Component is a function that returns HTML. Action is a function that returns HTML or calls another action. When you call action, it will be executed on server and result will be rendered in HTML document.

There are several ways to call actions:
- `ctx.Call(action, values...)` - call action with values
- `ctx.Submit(action)` - submit form with action
- `ctx.Click(action)` - click button with action
- `ctx.Send(action)` - send form with action

Depending on your needs you can render result of action in different places:
- `ctx.Replace(target)` - replace target with result
- `ctx.Render(target)` - render result inside target
- `ctx.None()` - do not render result

Code base is very simple, so please check examples and source code to see how it works.

## Components

### Basic Components

- `Div` - Container element
- `Span` - Inline text element
- `Form` - Form container
- `Input` - Input field
- `Button` - Button element
- `Select` - Dropdown select
- `Textarea` - Multi-line text input

### Input Types

- `IText` - Text input
- `IEmail` - Email input
- `IPhone` - Phone number input
- `INumber` - Numeric input
- `IPassword` - Password input

### Form Handling

```go
type LoginForm struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required"`
}

func (form *LoginForm) Login(ctx *ui.Context) string {
    if err := ctx.Body(form); err != nil {
        return form.Render(ctx, &err)
    }
    
    // Handle login logic
    return form.Success(ctx)
}
```

## Features

### Live Reload

Enable live reload during development:

```go
app.Autoreload()
```

### Session Management

```go
session := ctx.Session(db, "user")
session.Load(&userData)
session.Save(&userData)
```

### File Handling

```go
// Download file
ctx.DownloadAs(&fileReader, "application/pdf", "document.pdf")
```

### Form Validation

Built-in support for go-playground/validator:

```go
type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=18"`
}
```

## Styling

The framework integrates with Tailwind CSS by default. You can add custom styles through the `HtmlHead` method:

```go
app.HtmlHead = append(app.HtmlHead, `<link rel="stylesheet" href="custom.css">`)
```

## Examples

### Counter Component

A simple counter with increment and decrement buttons:

```go
type Counter struct {
    Count int
}

func (counter *Counter) Increment(ctx *ui.Context) string {
    ctx.Body(counter)
    counter.Count++
    return counter.Render(ctx)
}

func (counter *Counter) Decrement(ctx *ui.Context) string {
    ctx.Body(counter)
    counter.Count--
    if counter.Count < 0 {
        counter.Count = 0
    }
    return counter.Render(ctx)
}

func (counter *Counter) Render(ctx *ui.Context) string {
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
```

### Login Form with Validation

A complete login form with validation and error handling:

```go
type LoginForm struct {
    Name     string `validate:"required,oneof=user"`
    Password string `validate:"required,oneof=password"`
}

func (form *LoginForm) Login(ctx *ui.Context) string {
    if err := ctx.Body(form); err != nil {
        return form.Render(ctx, &err)
    }

    v := validator.New()
    if err := v.Struct(form); err != nil {
        return form.Render(ctx, &err)
    }

    return form.Success(ctx)
}

func (form *LoginForm) Success(ctx *ui.Context) string {
    return ui.Div("text-green-600 max-w-md p-8 text-center font-bold rounded-lg bg-white shadow-xl")("Success")
}

func (form *LoginForm) Render(ctx *ui.Context, err *error) string {
    target := ui.Target()
    login := ctx.Callable(form.Login)

    return ui.Form("flex flex-col gap-4 max-w-md bg-white p-8 rounded-lg shadow-xl", target, ctx.Submit(login).Replace(target))(
        ui.ErrorForm(err, &translations),

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
```

### Toggle Button

A button that toggles between two states:

```go
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

        return app.Html("Test", "p-8 bg-gray-200",
            ui.Div("flex flex-row gap-4")(
                ui.Div("flex justify-start gap-4 items-center")(
                    button(ctx),
                ),
            ),
        )
    }

    app.Page("/", page)
    app.Listen(":1422")
}
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 
