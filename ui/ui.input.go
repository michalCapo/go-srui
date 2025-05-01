package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type TInput struct {
	dates struct {
		Min time.Time
		Max time.Time
	}
	data         any
	Render       func(text string) string
	placeholder  string
	class        string
	classLabel   string
	classInput   string
	autocomplete string
	size         string
	onclick      string
	onchange     string
	as           string
	name         string
	pattern      string
	value        string
	error        validator.FieldError
	target       Attr
	numbers      struct {
		Min  float64
		Max  float64
		Step float64
	}
	visible  bool
	required bool
	disabled bool
}

func (c *TInput) Rows(value uint8) *TInput {
	c.target.Rows = value
	return c
}

func (c *TInput) If(value bool) *TInput {
	c.visible = value
	return c
}

func (c *TInput) Value(value string) *TInput {
	c.value = value
	return c
}

func (c *TInput) Type(value string) *TInput {
	c.as = value
	return c
}

func (c *TInput) Class(value ...string) *TInput {
	c.class = strings.Join(value, " ")
	return c
}

func (c *TInput) ClassInput(value ...string) *TInput {
	c.classInput = strings.Join(value, " ")
	return c
}

func (c *TInput) ClassLabel(value ...string) *TInput {
	c.classLabel = strings.Join(value, " ")
	return c
}

func (c *TInput) Size(value string) *TInput {
	c.size = value
	return c
}

func (c *TInput) Placeholder(value string) *TInput {
	c.placeholder = value
	return c
}

func (c *TInput) Pattern(value string) *TInput {
	c.pattern = value
	return c
}

func (c *TInput) Autocomplete(value string) *TInput {
	c.autocomplete = value
	return c
}

func (c *TInput) Required(value ...bool) *TInput {
	if value == nil {
		c.required = true
		return c
	}

	c.required = value[0]
	return c
}

func (c *TInput) Error(errs *error) *TInput {
	if errs == nil {
		return c
	}

	temp := (*errs).(validator.ValidationErrors)

	for _, err := range temp {
		if err.Field() == c.name {
			c.error = err
		}
	}

	return c
}

func (c *TInput) Disabled(value ...bool) *TInput {
	if value == nil {
		c.disabled = true
		return c
	}

	c.disabled = value[0]
	return c
}

func (c *TInput) Change(action string) *TInput {
	// if action.Target.Id == "" {
	// 	action.Target = c.target
	// }

	c.onchange = action

	return c
}

func (c *TInput) Click(action string) *TInput {
	// if action.Target.Id == "" {
	// 	action.Target = c.target
	// }

	c.onclick = action

	return c
}

func (c *TInput) Numbers(min float64, max float64, step float64) *TInput {
	c.numbers.Min = min
	c.numbers.Max = max
	c.numbers.Step = step

	return c
}

func (c *TInput) Dates(min time.Time, max time.Time) *TInput {
	c.dates.Min = min
	c.dates.Max = max

	return c
}

func IText(name string, data ...any) *TInput {
	c := &TInput{
		as:      "text",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			// if tmp.IsValid() {
			// 	value = tmp.String()
			// }

			tmp, err := PathValue(c.data, c.name)
			if err == nil {
				value = fmt.Sprintf("%v", tmp.Interface())
			}
		}

		// if c.error != nil { c.class = Classes(c.class, "border-4 border-red-600") }

		return Div(c.class)(
			Label(&c.target).
				Class(c.classLabel).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" }), c.classInput),
				Attr{
					Id:           c.target.Id,
					Name:         c.name,
					Type:         c.as,
					OnChange:     c.onchange,
					OnClick:      c.onclick,
					Required:     c.required,
					Disabled:     c.disabled,
					Value:        value,
					Pattern:      c.pattern,
					Placeholder:  c.placeholder,
					Autocomplete: c.autocomplete,
				},
			),

			// ErrorField(c.error),
		)
	}

	return c
}

var IPhone = func(name string, data ...any) *TInput {
	return IText(name, data...).
		Type("tel").
		Autocomplete("tel").
		Placeholder("+421").
		Pattern("\\+[0-9]{10,14}")
}

var IEmail = func(name string, data ...any) *TInput {
	return IText(name, data...).
		Type("email").
		Autocomplete("email").
		Placeholder("name@gmail.com")
	// Pattern("^[a-z0-9._%+-]+@[a-z0-9-]+\\.[a-z0-9]{2,}$")
}

func IArea(name string, data ...any) *TInput {
	c := &TInput{
		as:      "text",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		value := ""

		if c.data != nil {
			tmp, err := PathValue(c.data, c.name)
			if err == nil {
				value = fmt.Sprintf("%v", tmp.Interface())
			}

			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			// if tmp.IsValid() {
			// 	value = tmp.String()
			// }
		}

		rows := uint8(5)

		if c.target.Rows > 0 {
			rows = uint8(c.target.Rows)
		}

		return Div(c.class)(
			Label(&c.target).
				Required(c.required).
				Render(text),

			Textarea(
				Classes(AREA, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Rows: rows,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			)(value),
		)
	}

	return c
}

func IPassword(name string, data ...any) *TInput {
	c := &TInput{
		as:      "password",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			// if tmp.IsValid() {
			// 	value = tmp.String()
			// }

			tmp, err := PathValue(c.data, c.name)
			if err == nil {
				value = fmt.Sprintf("%v", tmp.Interface())
			}
		}

		return Div("")(
			Label(&c.target).
				Class(c.classLabel).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, c.class, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Value: value,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			),
		)
	}

	return c
}

func IDate(name string, data ...any) *TInput {
	c := &TInput{
		as:      "date",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		min := ""
		max := ""
		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			tmp, err := PathValue(c.data, c.name)

			if err == nil {
				if timeValue, ok := tmp.Interface().(time.Time); ok {
					if !timeValue.IsZero() {
						value = timeValue.Format(time.DateOnly)
						// value = fmt.Sprintf("%04d-%02d-%02d", timeValue.Year(), timeValue.Month(), timeValue.Day())
					}
				}
			}
		}

		if !c.dates.Min.IsZero() {
			min = c.dates.Min.Format(time.DateOnly)
		}

		if !c.dates.Max.IsZero() {
			max = c.dates.Max.Format(time.DateOnly)
		}

		return Div(c.class)(
			Label(&c.target).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Min:   min,
					Max:   max,
					Value: value,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					OnChange:    c.onchange,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			),
		)
	}
	return c
}

func ITime(name string, data ...any) *TInput {
	c := &TInput{
		as:      "time",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		min := ""
		max := ""
		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			tmp, err := PathValue(c.data, c.name)

			if err == nil {
				if timeValue, ok := tmp.Interface().(time.Time); ok {
					value = timeValue.Format("15:04")
					// value = fmt.Sprintf("%02d:%02d", timeValue.Hour(), timeValue.Minute())
				}
			}
		}

		if !c.dates.Min.IsZero() {
			min = c.dates.Min.Format("15:04")
		}

		if !c.dates.Max.IsZero() {
			max = c.dates.Max.Format("15:04")
		}

		return Div("")(
			Label(&c.target).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, c.class, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Min:   min,
					Max:   max,
					Value: value,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			),
		)
	}
	return c
}

func IDateTime(name string, data ...any) *TInput {
	c := &TInput{
		as:      "datetime-local",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		min := ""
		max := ""
		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			tmp, err := PathValue(c.data, c.name)

			if err == nil {
				if timeValue, ok := tmp.Interface().(time.Time); ok {
					// value = fmt.Sprintf("%04d-%02d-%02d %02d:%02d", timeValue.Year(), timeValue.Month(), timeValue.Day(), timeValue.Hour(), timeValue.Minute())
					value = timeValue.Format("2006-01-02T15:04")
				}
			}
		}

		if !c.dates.Min.IsZero() {
			min = c.dates.Min.Format("2006-01-02T15:04")
		}

		if !c.dates.Max.IsZero() {
			max = c.dates.Max.Format("2006-01-02T15:04")
		}

		return Div("")(
			Label(&c.target).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, c.class, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Min:   min,
					Max:   max,
					Value: value,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			),
		)
	}
	return c
}

func INumber(name string, data ...any) *TInput {
	c := &TInput{
		as:      "number",
		target:  Target(),
		name:    name,
		size:    MD,
		visible: true,
	}

	if len(data) > 0 {
		c.data = data[0]
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		min := ""
		max := ""
		step := ""
		value := ""

		if c.data != nil {
			// v := reflect.ValueOf(c.data)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			// tmp := v.FieldByName(c.name)

			// if tmp.IsValid() {
			// 	value = fmt.Sprintf("%v", tmp.Interface())
			// }

			tmp, err := PathValue(c.data, c.name)
			if err == nil {
				value = fmt.Sprintf("%v", tmp.Interface())
			}
		}

		if c.numbers.Min != 0 {
			min = fmt.Sprintf("%v", c.numbers.Min)
		}

		if c.numbers.Max != 0 {
			max = fmt.Sprintf("%v", c.numbers.Max)
		}

		if c.numbers.Step != 0 {
			step = fmt.Sprintf("%v", c.numbers.Step)
		}

		return Div(c.class)(
			Label(&c.target).
				Required(c.required).
				Render(text),

			Input(
				Classes(INPUT, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
				Attr{
					Min:   min,
					Max:   max,
					Step:  step,
					Value: value,

					Type:        c.as,
					Id:          c.target.Id,
					Name:        c.name,
					OnClick:     c.onclick,
					Required:    c.required,
					Disabled:    c.disabled,
					Placeholder: c.placeholder,
				},
			),
		)
	}
	return c
}

var Hidden = func(name string, typ string, value any, attr ...Attr) string {
	return Input("hidden", append(attr, Attr{Name: name, Type: typ, Value: fmt.Sprintf("%v", value)})...)
}

// func IValue(name string, data ...any) *TInput {
// 	c := &TInput{
// 		as:      "text",
// 		target:  Target(),
// 		name:    name,
// 		size:    MD,
// 		visible: true,
// 		data:    data[0],
// 	}

// 	c.Render = func(text string) string {
// 		if !c.visible {
// 			return ""
// 		}

// 		value := ""

// 		if c.data != nil {
// 			tmp, err := PathValue(c.data, c.name)
// 			if err == nil {
// 				value = fmt.Sprintf("%v", tmp.Interface())
// 			}
// 		}

// 		return Div(c.class)(
// 			Label(c.target).
// 				Required(c.required).
// 				Render(text),

// 			Div(
// 				Classes(VALUE, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" })),
// 				Attr{
// 					Id:           c.target.Id,
// 					Name:         c.name,
// 					Type:         c.as,
// 					OnChange:     c.onchange,
// 					OnClick:      c.onclick,
// 					Required:     c.required,
// 					Disabled:     c.disabled,
// 					Value:        value,
// 					Pattern:      c.pattern,
// 					Placeholder:  c.placeholder,
// 					Autocomplete: c.autocomplete,
// 				},
// 			)(),
// 		)
// 	}

// 	return c
// }

func IValue(attr ...Attr) *TInput {
	c := &TInput{
		target:  Target(),
		size:    MD,
		visible: true,
	}

	c.Render = func(text string) string {
		if !c.visible {
			return ""
		}

		attr = append(attr, Attr{
			Id:          c.target.Id,
			Name:        c.name,
			Required:    c.required,
			Disabled:    c.disabled,
			Pattern:     c.pattern,
			Placeholder: c.placeholder,
		})

		return Div(c.class)(
			Label(&c.target).
				Class(c.classLabel).
				Required(c.required).
				Render(text),

			Div(
				Classes(VALUE, c.size, If(c.disabled, func() string { return DISABLED }), If(c.error != nil, func() string { return "border-l-8 border-red-600" }), c.classInput),
				attr...,
			)(c.value),
		)
	}

	return c
}
