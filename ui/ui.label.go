package ui

type ALabel struct {
	id         string
	class      string
	classLabel string
	required   bool
	disabled   bool
}

func Label(target *Attr) *ALabel {
	tmp := &ALabel{
		class: "text-sm",
	}

	if target != nil {
		tmp.id = target.Id
	}

	return tmp
}

func (c *ALabel) Required(value bool) *ALabel {
	c.required = value
	return c
}

func (c *ALabel) Disabled(value bool) *ALabel {
	c.disabled = value
	return c
}

func (c *ALabel) Class(value ...string) *ALabel {
	c.class = Classes(value...)
	return c
}

func (c *ALabel) ClassLabel(value ...string) *ALabel {
	c.classLabel = Classes(value...)
	return c
}

func (c *ALabel) Render(text string) string {
	if text == "" {
		return ""
	}

	return Div(Classes(c.class, "relative"))(
		open("label")(c.classLabel, Attr{For: c.id})(text),
		If(c.required && !c.disabled, func() string { return Span("ml-1 text-red-700")("*") }),
	)
}
