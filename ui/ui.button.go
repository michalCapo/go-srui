package ui

import "strings"

type AButton struct {
	size     string
	color    string
	onclick  string
	class    string
	as       string
	target   Attr
	visible  bool
	disabled bool
	attr     []Attr
}

func Button(attr ...Attr) *AButton {
	return &AButton{
		as:      "div",
		size:    MD,
		visible: true,
		target:  Target(),
		attr:    attr,
	}
}

func (b *AButton) Submit() *AButton {
	b.as = "submit"
	return b
}

func (b *AButton) Reset() *AButton {
	b.as = "reset"
	return b
}

func (b *AButton) If(value bool) *AButton {
	b.visible = value
	return b
}

func (b *AButton) Disabled(value bool) *AButton {
	b.disabled = value
	return b
}

func (b *AButton) Class(value ...string) *AButton {
	b.class = strings.Join(value, " ")
	return b
}

func (b *AButton) Color(value string) *AButton {
	b.color = value
	return b
}

func (b *AButton) Size(value string) *AButton {
	b.size = value
	return b
}

func (b *AButton) Click(onclick string) *AButton {
	// if action.Target.Id == "" {
	// 	action.Target = b.target
	// }

	b.onclick = onclick

	return b
}

func (b *AButton) Href(value string) *AButton {
	b.as = "a"
	b.attr = append(b.attr, Attr{Href: value})
	return b
}

func (b *AButton) Render(text string) string {
	if !b.visible {
		return ""
	}

	class := Classes(BTN, b.size, b.color, b.class, If(b.disabled, func() string { return DISABLED + " opacity-25" }))

	if b.as == "a" {
		return A(
			class,
			append(b.attr, Attr{
				Id: b.target.Id,
			})...,
		)(text)
	}

	if b.as == "div" {
		return Div(
			class,
			append(b.attr, Attr{
				Id:      b.target.Id,
				OnClick: b.onclick,
			})...,
		)(text)
	}

	return open("button")(
		class,
		append(b.attr,
			Attr{
				Id:      b.target.Id,
				Type:    b.as,
				OnClick: b.onclick,
			})...,
	)(text)
}
