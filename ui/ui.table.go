package ui

import (
	"fmt"
	"strings"
)

type TTable[T any] struct {
	heads []string
	slots []struct {
		slot func(item *T) string
		cls  string
	}
	data []*T
	cls  string
	// rows []string
}

func Table[T any](cls string, data []*T) *TTable[T] {
	return &TTable[T]{
		heads: []string{},
		slots: []struct {
			slot func(item *T) string
			cls  string
		}{},
		data: data,
		cls:  cls,
	}
}

func NewTable[T any](cls string) *TTable[T] {
	return &TTable[T]{}
}

func (t *TTable[T]) Head(value string, cls string) *TTable[T] {
	t.heads = append(t.heads, fmt.Sprintf(`<th class="%s">%s</th>`, cls, value))
	return t
}

func (t *TTable[T]) Field(slot func(item *T) string, cls string) *TTable[T] {
	t.slots = append(t.slots, struct {
		slot func(item *T) string
		cls  string
	}{slot, cls})
	return t
}

func (t *TTable[T]) Row(slot func(item *T) []string, cls string) *TTable[T] {
	return t
}

func (t *TTable[T]) Render() string {
	var headsBuilder strings.Builder
	for _, head := range t.heads {
		headsBuilder.WriteString(head)
	}

	var rowsBuilder strings.Builder
	for _, row := range t.data {
		rowsBuilder.WriteString("<tr>")
		for _, slot := range t.slots {
			rowsBuilder.WriteString(fmt.Sprintf(`<td class="%s">%s</td>`, slot.cls, slot.slot(row)))
		}
		rowsBuilder.WriteString("</tr>")
	}

	return fmt.Sprintf(
		`<div><table class="table-auto %s"><thead><tr>%s</tr></thead><tbody>%s</tbody></table></div>`,
		t.cls, headsBuilder.String(), rowsBuilder.String(),
	)
}
