package ui

// func Times(class ...string) string {
// 	return Div(Classes(class...))("&#215;")
// }

// func Plus(class ...string) string {
// 	return Div(Classes(class...))("&#43;")
// }

// func Minus(class ...string) string {
// 	return Div(Classes(class...))("&#8722;")
// }

// func Check(class ...string) string {
// 	return Div(Classes(class...))("&#10003;")
// }

// func Close(class ...string) string {
// 	return Div(Classes(class...))("&#10005;")
// }

// func ArrowUp(class ...string) string {
// 	return Div(Classes(class...))("&#8593;")
// }

// func ArrowDown(class ...string) string {
// 	return Div(Classes(class...))("&#8595;")
// }

// func ArrowLeft(class ...string) string {
// 	return Div(Classes(class...))("&#8592;")
// }

// func ArrowRight(class ...string) string {
// 	return Div(Classes(class...))("&#8594;")
// }

func Icon(class string, attr ...Attr) string {
	return Div(class, attr...)()
}

func Icon2(class string, text string) string {
	return Div("flex-1 flex items-center gap-2")(
		Icon(class),
		Flex1,
		Div("text-center")(text),
		Flex1,
	)
}
