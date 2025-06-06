package ui

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

const (
	XS = " p-1"
	SM = " p-2"
	MD = " p-3"
	ST = " p-4"
	LG = " p-5"
	XL = " p-6"
)

const (
	AREA          = " cursor-pointer bg-white border border-gray-400 hover:border-blue-500 rounded block w-full shadow"
	INPUT         = " cursor-pointer bg-white border border-gray-400 hover:border-blue-500 rounded block w-full shadow h-12"
	VALUE         = " bg-white border border-gray-400 hover:border-blue-500 rounded block shadow h-12"
	BTN           = " cursor-pointer font-bold text-center select-none"
	DISABLED      = " cursor-text bg-gray-300 pointer-events-none"
	Yellow        = " bg-yellow-400 text-gray-800 hover:text-gray-200 hover:bg-yellow-600 font-bold border-gray-300 flex items-center justify-center"
	YellowOutline = " border border-yellow-500 text-yellow-600 hover:text-gray-700 hover:bg-yellow-500 flex items-center justify-center"
	Green         = " bg-green-600 text-white hover:bg-green-700 checked:bg-green-600 border-gray-300 flex items-center justify-center"
	GreenOutline  = " border border-green-500 text-green-500 hover:text-white hover:bg-green-599 flex items-center justify-center"
	Purple        = " bg-purple-500 text-white hover:bg-purple-700 border-purple-500 flex items-center justify-center"
	PurpleOutline = " border border-purple-500 text-purple-500 hover:text-white hover:bg-purple-600 flex items-center justify-center"
	Blue          = " bg-blue-800 text-white hover:bg-blue-700 border-gray-300 flex items-center justify-center"
	BlueOutline   = " border border-blue-500 text-blue-600 hover:text-white hover:bg-blue-700 checked:bg-blue-700 flex items-center justify-center"
	Red           = " bg-red-600 text-white hover:bg-red-800 border-gray-300 flex items-center justify-center"
	RedOutline    = " border border-red-500 text-red-600 hover:text-white hover:bg-red-700 flex items-center justify-center"
	Gray          = " bg-gray-600 text-white hover:bg-gray-800 focus:bg-gray-800 border-gray-300 flex items-center justify-center"
	GrayOutline   = " border border-gray-300 text-black hover:text-white hover:bg-gray-700 flex items-center justify-center"
	White         = " bg-white text-black hover:bg-gray-200 border-gray-200 flex items-center justify-center"
	WhiteOutline  = " border border-white text-balck hover:text-black hover:bg-white flex items-center justify-center"
)

type Attr struct {
	OnClick      string
	Step         string
	Id           string
	Href         string
	Title        string
	Alt          string
	Type         string
	Class        string
	Style        string
	Name         string
	Value        string
	Checked      string
	OnSubmit     string
	For          string
	Src          string
	Selected     string
	Pattern      string
	Placeholder  string
	Autocomplete string
	OnChange     string
	Max          string
	Min          string
	Target       string
	Rows         uint8
	Cols         uint8
	Disabled     bool
	Required     bool
}

type AOption struct {
	Id    string
	Value string
}

func MakeOptions(options []string) []AOption {
	var result []AOption

	for _, option := range options {
		result = append(result, AOption{Id: option, Value: option})
	}

	return result
}

func attributes(attrs ...Attr) string {
	var result []string

	for _, attr := range attrs {

		if attr.Id != "" {
			result = append(result, fmt.Sprintf(`id="%s"`, attr.Id))
		}

		if attr.Href != "" {
			result = append(result, fmt.Sprintf(`href="%s"`, attr.Href))
		}

		if attr.Alt != "" {
			result = append(result, fmt.Sprintf(`alt="%s"`, attr.Alt))
		}

		if attr.Title != "" {
			result = append(result, fmt.Sprintf(`title="%s"`, attr.Title))
		}

		if attr.Src != "" {
			result = append(result, fmt.Sprintf(`src="%s"`, attr.Src))
		}

		if attr.For != "" {
			result = append(result, fmt.Sprintf(`for="%s"`, attr.For))
		}

		if attr.Type != "" {
			result = append(result, fmt.Sprintf(`type="%s"`, attr.Type))
		}

		if attr.Class != "" {
			result = append(result, fmt.Sprintf(`class="%s"`, attr.Class))
		}

		if attr.Style != "" {
			result = append(result, fmt.Sprintf(`style="%s"`, attr.Style))
		}

		if attr.OnClick != "" {
			result = append(result, fmt.Sprintf(`onclick="%s"`, attr.OnClick))
		}

		if attr.OnChange != "" {
			result = append(result, fmt.Sprintf(`onchange="%s"`, attr.OnChange))
		}

		if attr.OnSubmit != "" {
			result = append(result, fmt.Sprintf(`onsubmit="%s"`, attr.OnSubmit))
		}

		if attr.Value != "" {
			result = append(result, fmt.Sprintf(`value="%s"`, attr.Value))
		}

		if attr.Checked != "" {
			result = append(result, fmt.Sprintf(`checked="%s"`, attr.Checked))
		}

		if attr.Selected != "" {
			result = append(result, fmt.Sprintf(`selected="%s"`, attr.Selected))
		}

		if attr.Name != "" {
			result = append(result, fmt.Sprintf(`name="%s"`, attr.Name))
		}

		if attr.Placeholder != "" {
			result = append(result, fmt.Sprintf(`placeholder="%s"`, attr.Placeholder))
		}

		if attr.Autocomplete != "" {
			result = append(result, fmt.Sprintf(`autocomplete="%s"`, attr.Autocomplete))
		}

		if attr.Pattern != "" {
			result = append(result, fmt.Sprintf(`pattern="%s"`, attr.Pattern))
		}

		if attr.Cols != 0 {
			result = append(result, fmt.Sprintf(`cols="%d"`, attr.Cols))
		}

		if attr.Rows != 0 {
			result = append(result, fmt.Sprintf(`rows="%d"`, attr.Rows))
		}

		if attr.Step != "" {
			result = append(result, fmt.Sprintf(`step="%s"`, attr.Step))
		}

		if attr.Min != "" {
			result = append(result, fmt.Sprintf(`min="%s"`, attr.Min))
		}

		if attr.Max != "" {
			result = append(result, fmt.Sprintf(`max="%s"`, attr.Max))
		}

		if attr.Target != "" {
			result = append(result, fmt.Sprintf(`target="%s"`, attr.Target))
		}

		if attr.Required {
			result = append(result, `required="required"`)
		}

		if attr.Disabled {
			result = append(result, `disabled="disabled"`)
		}
	}

	return strings.Join(result, " ")
}

func open(tag string) func(class string, attr ...Attr) func(elements ...string) string {
	return func(class string, attr ...Attr) func(elements ...string) string {
		return func(elements ...string) string {
			attr = append(attr, Attr{Class: class})
			return fmt.Sprintf(`<%s %s>%s</%s>`, tag, attributes(attr...), strings.Join(elements, " "), tag)
		}
	}
}

func closed(tag string) func(class string, attr ...Attr) string {
	return func(class string, attr ...Attr) string {
		attr = append(attr, Attr{Class: class})
		return fmt.Sprintf(`<%s %s/>`, tag, attributes(attr...))
	}
}

var (
	I        = open("i")
	A        = open("a")
	Div      = open("div")
	Span     = open("span")
	Form     = open("form")
	Textarea = open("textarea")
	Select   = open("select")
	Option   = open("option")
	List     = open("ul")
	ListItem = open("li")

	Img   = closed("img")
	Input = closed("input")

	Flex1 = Div("flex-1")()

	Space = "&nbsp;"
	W35   = Attr{Style: "width: 35rem;"}
	W30   = Attr{Style: "width: 30rem;"}
	W25   = Attr{Style: "width: 25rem;"}
	W20   = Attr{Style: "width: 20rem;"}
)

var Script = func(value ...string) string {
	return Trim(fmt.Sprintf(`<script>%s</script>`, strings.Join(value, " ")))
}

var Target = func() Attr {
	return Attr{Id: "i" + RandomString(15)}
}

func Variable[T any](getter func(*T) string) func(item *T) Attr {
	temp := Target()

	return func(item *T) Attr {
		value := getter(item)
		return Attr{Id: temp.Id + "_" + value}
	}
}

var Id = func(target string) Attr {
	return Attr{Id: target}
}

var Href = func(value string, target ...string) Attr {
	if len(target) > 0 {
		return Attr{Href: value, Target: target[0]}
	}

	return Attr{Href: value}
}

var Src = func(alt string, target string) Attr {
	return Attr{Src: target, Alt: alt}
}

var Title = func(value string) Attr {
	return Attr{Title: string(value)}
}

var Classes = func(values ...string) string {
	return Trim(strings.Join(values, " "))
}

// type Hook struct {
// 	Target string
// }

// func (h *Hook) Id() Attr {
// 	if len(h.Target) == 0 {
// 		h.Target = "i" + RandomString(15)
// 	}

// 	return Attr{Id: h.Target}
// }

// func (h *Hook) Submit(ctx *Context, action Action) Attr {
// 	action.Type = FORM

// 	if action.Target.Id == "" {
// 		action.Target = Attr{Id: h.Target}
// 	}

// 	return Attr{OnSubmit: ctx.Post(action)}
// }

// func (h *Hook) Click(ctx *Context, action Action) Attr {
// 	if action.Target.Id == "" {
// 		action.Target = Attr{Id: h.Target}
// 	}

// 	return Attr{OnClick: ctx.Post(action)}
// }

type Action struct {
	Method *Callable
	Target Attr
	Values []any
	// Type   ActionType
}

// func SendForm(Target *Attr, Method ComponentMethod, Values ...any) Action {
// 	return Action{
// 		Type:   FORM,
// 		Target: *Target,
// 		Method: Method,
// 		Values: Values,
// 	}
// }

var (
	re  = regexp.MustCompile(`\s{4,}`)
	re2 = regexp.MustCompile(`[\t\n]+`)
	re3 = regexp.MustCompile(`"`)
)

func Trim(s string) string {
	return re.ReplaceAllString(re2.ReplaceAllString(s, ""), " ")
}

func Normalize(s string) string {
	return re.ReplaceAllString(re2.ReplaceAllString(re3.ReplaceAllString(s, "&quot;"), ""), " ")
	// return re3.ReplaceAllString(s, "&quot;")
}

func If(cond bool, value func() string) string {
	if cond {
		return value()
	}

	return ""
}

func Iff(cond bool) func(value ...string) string {
	return func(value ...string) string {
		if cond {
			return strings.Join(value, " ")
		}

		return ""
	}
}

func Or(cond bool, value func() string, other func() string) string {
	if cond {
		return value()
	}

	return other()
}

// func Orr(cond bool) (func(...string) string, func(...string) string) {
// 	yes := func(value ...string) string {
// 		if cond {
// 			return strings.Join(value, " ")
// 		}

// 		return ""
// 	}

// 	no := func(value ...string) string {
// 		if !cond {
// 			return strings.Join(value, " ")
// 		}

// 		return ""
// 	}

// 	return yes, no
// }

func Map[T any](values []T, iter func(*T, int) string) string {
	var result []string
	for key, value := range values {
		result = append(result, iter(&value, key))
	}

	return strings.Join(result, " ")
}

func Map2[T any](values []T, iter func(T, int) []string) string {
	var result []string
	for key, value := range values {
		result = append(result, iter(value, key)...)
	}

	return strings.Join(result, " ")
}
func For[T any](from, to int, iter func(int) string) string {
	var result []string
	for i := from; i < to; i++ {
		result = append(result, iter(i))
	}

	return strings.Join(result, " ")
}

func Error(err error) string {
	stackErr := errors.WithStack(err)
	fmt.Printf("%+v\n", stackErr)

	return Div("bg-red-500 text-white font-bold p-8 text-center border border-red-800")(
		"Opps, something went wrong",
	)
}

func ErrorField(err validator.FieldError) string {
	if err == nil {
		return ""
	}

	return Div("text-red-600 p-2 rounded bg-white")(
		// Span("bg-red-800 text-white rounded-full")(err.Field()),
		"obsahuje nevalidnú hodnotu",
	)
}

func ErrorForm(errs *error, translations *map[string]string) string {
	if errs == nil || *errs == nil {
		return ""
	}

	temp := (*errs).(validator.ValidationErrors)

	// for _, err := range errs { fmt.Printf("%+v\n", err.Field()) }

	return Div("text-red-600 p-4 rounded text-center border-4 border-red-600 bg-white")(
		Map(temp, func(err *validator.FieldError, _ int) string {
			trans := (*err).Field()
			invalid := "has invalid value"

			if translations != nil && (*translations)[trans] != "" {
				trans = (*translations)[trans]
			}

			if translations != nil && (*translations)[invalid] != "" {
				invalid = (*translations)[invalid]
			}

			return Div("")(
				Span("font-bold uppercase")(trans),
				Space,
				invalid,
			)
		}),
	)
}

func Print(value any) string {
	return fmt.Sprintf("%+v", value)
}

func PathValue(obj any, path string) (*reflect.Value, error) {
	parts := strings.Split(path, ".")
	current := reflect.ValueOf(obj)

	for _, part := range parts {
		// Handle slice index notation
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			fieldName := part[:strings.Index(part, "[")]
			indexStr := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]
			indexVal, err := strconv.Atoi(indexStr)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			current = current.Elem().FieldByName(fieldName)
			if !current.IsValid() || current.Kind() != reflect.Slice {
				fmt.Printf("invalid slice field: %s", fieldName)
				return nil, err
			}

			for current.Len() <= indexVal {
				elemType := current.Type().Elem()

				if elemType.Kind() == reflect.Ptr {
					newElem := reflect.New(elemType.Elem())
					current.Set(reflect.Append(current, newElem))
				} else {
					newElem := reflect.New(elemType).Elem()
					current.Set(reflect.Append(current, newElem))
				}
			}

			current = current.Index(indexVal)

			// fmt.Printf("current: %v\n", current)
		} else {
			if current.Kind() == reflect.Ptr {
				current = current.Elem()
			}

			current = current.FieldByName(part)
		}

		if !current.IsValid() {
			err := fmt.Errorf("invalid path segment: %s", part)

			fmt.Println(err)

			return nil, err
		}
	}

	return &current, nil
}
