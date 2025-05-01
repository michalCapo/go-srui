package ui

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

func Icon3(class string, text string) string {
	return Div("flex-1 flex items-center gap-2")(
		Flex1,
		Div("text-center")(text),
		Icon(class),
		Flex1,
	)
}

func Icon4(class string, text string) string {
	return Div("flex-1 flex items-center gap-2")(
		Flex1,
		Div("text-center")(text),
		Flex1,
		Icon(class),
	)
}
