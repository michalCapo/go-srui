package ui

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// type TSort struct {
// 	Field string
// 	Text  string
// }

const (
	BOOL = iota
	NOT_ZERO_DATE
	DATES
)

type TField struct {
	DB    string
	Field string
	Text  string
	As    uint
	Bool  bool
	Dates struct {
		From time.Time
		To   time.Time
	}
}

type TQuery struct {
	Limit  int64
	Offset int64
	Order  string
	Search string
	Filter []TField
}

type TCollateResult[T any] struct {
	Total    int64
	Filtered int64
	Data     []T
	Query    *TQuery
}

type TCollate[T any] struct {
	Limit        int64
	Target       Attr
	TargetFilter Attr
	Database     *gorm.DB
	onReset      **Callable
	onResize     **Callable
	onSort       **Callable
	onSearch     **Callable
	onXLS        **Callable
	Search       []TField
	Sort         []TField
	Filter       []TField
	Excel        []TField
	OnRow        func(*T, int) string
	OnExcel      func(*[]T) (string, io.Reader, error)
	Set          func(*TQuery)
	Get          func(*TQuery)
}

type CollateMethod[T any] = func(*Context, *TCollate[T])

func Collate[T any](app *App, prefix string, init *TQuery, setup CollateMethod[T]) func(*Context, *TSession, *gorm.DB) func() string {

	collate := &TCollate[T]{
		Target:       Target(),
		TargetFilter: Target(),
	}

	collate.onXLS = app.Action(prefix+"-collate-xls", func(ctx *Context) string {
		// Set query for all records
		query := makeQuery(init)

		if collate.Get != nil {
			collate.Get(query)
		}

		query.Limit = 1000000
		result := collate.Load(query)

		var filename string
		var reader io.Reader

		if collate.OnExcel != nil {
			var err error

			filename, reader, err = collate.OnExcel(&result.Data)
			if err != nil {
				fmt.Println(err)
				return "Error generating Excel file"
			}
		} else {
			f := excelize.NewFile()
			defer f.Close()

			for i, header := range collate.Excel {
				if header.Text == "" {
					header.Text = header.Field
				}

				cell := string(rune('A'+i)) + "1"
				f.SetCellValue("Sheet1", cell, header.Text)
			}

			styleDate, err := f.NewStyle(&excelize.Style{NumFmt: 14})
			if err != nil {
				fmt.Println(err)
			}

			// Write data rows
			for rowIndex, item := range result.Data {
				v := reflect.ValueOf(item)

				for colIndex, header := range collate.Excel {
					col := string(rune('A' + colIndex))
					cell := col + strconv.Itoa(rowIndex+2)
					value := v.FieldByName(header.Field).Interface()
					typ := v.FieldByName(header.Field).Type().String()

					switch typ {
					case "time.Time":
						if !value.(time.Time).IsZero() {
							// value = value.(time.Time).Format("2006-01-02")

							f.SetCellValue("Sheet1", cell, value)
							f.SetCellStyle("Sheet1", cell, cell, styleDate)

							f.SetColWidth("Sheet1", col, col, 15)
						}

					default:
						f.SetCellValue("Sheet1", cell, value)
					}
				}
			}

			// Set filename with timestamp
			filename = fmt.Sprintf("export_%s.xlsx", time.Now().Format("20060102_150405"))

			fileBytes, err := f.WriteToBuffer()
			if err != nil {
				return "Error generating Excel file"
			}

			reader = io.Reader(bytes.NewReader(fileBytes.Bytes()))
		}

		ctx.DownloadAs(&reader, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", filename)

		return ""
	})

	collate.onResize = app.Action(prefix+"-collate-resize", func(ctx *Context) string {
		query := makeQuery(init)

		if collate.Get != nil {
			collate.Get(query)
		}

		query.Limit += query.Limit

		if collate.Set != nil {
			collate.Set(query)
		}

		return collate.Render(ctx, query)
	})

	collate.onSort = app.Action(prefix+"-collate-sort", func(ctx *Context) string {
		query := makeQuery(init)

		if collate.Get != nil {
			collate.Get(query)
		}

		body := &TQuery{}
		err := ctx.Body(body)
		if err != nil {
			fmt.Println(err)
		}

		query.Order = body.Order

		collate.Set(query)

		return collate.Render(ctx, query)
	})

	collate.onSearch = app.Action(prefix+"-collate-search", func(ctx *Context) string {
		query := makeQuery(init)

		if collate.Get != nil {
			collate.Get(query)
		}

		body := &TQuery{}
		err := ctx.Body(body)
		if err != nil {
			fmt.Println(err)
		}

		query.Search = body.Search
		query.Filter = body.Filter

		if collate.Set != nil {
			collate.Set(query)
		}

		return collate.Render(ctx, query)
	})

	collate.onReset = app.Action(prefix+"-collate-reset", func(ctx *Context) string {
		query := makeQuery(init)

		// if collate.Get != nil {
		// 	collate.Get(query)
		// }

		query.Limit = collate.Limit

		if collate.Set != nil {
			collate.Set(query)
		}

		return collate.Render(ctx, query)
	})

	return func(ctx *Context, session *TSession, database *gorm.DB) func() string {
		collate.Database = database
		collate.Set = func(query *TQuery) {
			session.Save(query)
		}
		collate.Get = func(query *TQuery) {
			session.Load(query)
		}

		setup(ctx, collate)

		query := makeQuery(init)

		if collate.Get != nil {
			collate.Get(query)
		}

		if collate.Set != nil {
			collate.Set(query)
		}

		return func() string {
			return collate.Render(ctx, query)
		}
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func (c *TCollate[T]) Load(query *TQuery) *TCollateResult[T] {
	result := &TCollateResult[T]{
		Total:    0,
		Filtered: 0,
		Data:     []T{},
		Query:    query,
	}

	c.Database.Model(result.Data).Count(&result.Total)

	temp := c.Database.Model(&result.Data).
		Session(&gorm.Session{}).
		Order(query.Order).
		Limit(int(query.Limit)).
		Offset(int(query.Offset))

	var sorts []string
	var filters []string

	for _, filter := range query.Filter {
		if filter.DB == "" {
			filter.DB = filter.Field
		}

		if filter.As == NOT_ZERO_DATE && filter.Bool {
			// temp.Where(filter.Field + " > '0001-01-01 00:00:00+00:00'")
			filters = append(filters, filter.DB+" > '0001-01-01 00:00:00+00:00'")
		}

		if filter.As == DATES && !filter.Dates.From.IsZero() {
			// temp.Where(filter.Field+" BETWEEN ? AND ?", filter.Dates.From, filter.Dates.To)
			filters = append(filters, filter.DB+" >= '"+startOfDay(filter.Dates.From).Format("2006-01-02 15:04:05")+"'")
		}

		if filter.As == DATES && !filter.Dates.To.IsZero() {
			// temp.Where(filter.Field+" BETWEEN ? AND ?", filter.Dates.Frobm, filter.Dates.To)

			filters = append(filters, filter.DB+" <= '"+endOfDay(filter.Dates.To).Format("2006-01-02 15:04:05")+"'")
		}
	}

	if len(filters) > 0 {
		temp = temp.Where("(" + strings.Join(filters, " AND ") + ")")
	}

	if len(query.Search) > 0 {
		for _, field := range c.Search {
			if field.DB == "" {
				field.DB = field.Field
			}

			// temp.Or(field+" LIKE ?", "%"+query.Search+"%")
			sorts = append(sorts, field.DB+" LIKE '%"+query.Search+"%'")
		}
	}

	if len(sorts) > 0 {
		temp = temp.Where("(" + strings.Join(sorts, " OR ") + ")")
	}

	temp.Count(&result.Filtered)
	temp.Find(&result.Data)

	return result
}

func (collate *TCollate[T]) Render(ctx *Context, query *TQuery) string {
	result := collate.Load(query)

	return Div("flex flex-col gap-2", collate.Target)(
		Div("flex gap-x-2")(
			Sorting(ctx, collate, query),
			Flex1,
			Searching(ctx, collate, query),
		),
		Div("flex justify-end")(
			Filtering(ctx, collate, query),
		),
		Map(result.Data, collate.OnRow),
		Paging(ctx, collate, result),
	)
}

func makeQuery(def *TQuery) *TQuery {
	if def.Offset < 0 {
		def.Offset = 0
	}

	if def.Limit <= 0 {
		def.Limit = 10
	}

	query := &TQuery{
		Limit:  def.Limit,
		Offset: def.Offset,
		Order:  def.Order,
		Search: def.Search,
		Filter: def.Filter,
	}

	return query
}

func Empty[T any](result *TCollateResult[T]) string {
	if result.Total == 0 {
		return Div("mt-2 py-24 rounded text-xl flex justify-center items-center bg-white rounded-lg")(
			Div("")(
				Div("text-black text-2xl p-4 mb-2 font-bold flex justify-center items-center")(("Nenašli sa žiadne záznamy")),
			),
		)
	}

	if result.Filtered == 0 {
		return Div("mt-2 py-24 rounded text-xl flex justify-center items-center bg-white rounded-lg")(
			Div("flex gap-x-px items-center justify-center text-2xl")(
				Icon("fa fa-fw fa-exclamation-triangle text-yellow-500"),
				Div("text-black p-4 mb-2 font-bold flex justify-center items-center")("Nenašli sa žiadne záznamy pre zvolený filter"),
			),
		)
	}

	return ""
}

func Filtering[T any](ctx *Context, collate *TCollate[T], query *TQuery) string {
	if len(collate.Filter) == 0 {
		return ""
	}

	// c.Query = Query(def)
	// ctx.Session(database, name, c.Query)

	return Div("col-span-2 relative h-0 hidden z-20", collate.TargetFilter)(
		Div("absolute top-1 right-0 rounded-lg bg-gray-100 border border-black shadow-2xl p-4")(
			Form("flex flex-col", ctx.Send(collate.onSearch).AsSubmit(collate.Target))(
				Hidden("Search", "string", query.Search),

				Map2(collate.Filter, func(item TField, index int) []string {
					if item.DB == "" {
						item.DB = item.Field
					}

					position := fmt.Sprintf("Filter[%d]", index)

					return []string{
						Iff(item.As == NOT_ZERO_DATE)(
							Div("flex")(
								Hidden(position+".Field", "string", item.DB),
								Hidden(position+".As", "uint", item.As),
								ICheckbox(position+".Bool", query).
									Render(item.Text),
							),
						),

						Iff(item.As == DATES)(
							Label(nil).Class("text-xs mt-3 font-bold").Render(item.Text),
							Div("flex gap-1")(
								Hidden(position+".Field", "string", item.DB),
								Hidden(position+".As", "uint", item.As),
								IDate(position+".Dates.From", query).Render(""),
								IDate(position+".Dates.To", query).Render(""),
							),
						),
					}
				}),

				Div("flex gap-px mt-3")(
					Button().
						Color(Red).
						Click(ctx.Call(collate.onReset).Replace(collate.Target)).
						Class("rounded-l-lg").
						Render(Icon("fa fa-fw fa-times")),

					Button().
						Submit().
						Class("flex-1 rounded-r-lg").
						Color(Purple).
						Render(Icon2("fa fa-fw fa-search", "Filtrovať")),
				),
			),
		),
	)
}

func Searching[T any](ctx *Context, collate *TCollate[T], query *TQuery) string {
	if collate.Search == nil {
		return ""
	}

	// reset := TQuery{
	// 	Search: "",
	// 	Filter: query.Filter,
	// 	Order:  query.Order,
	// 	Offset: query.Offset,
	// 	Limit:  query.Limit,
	// }

	return Div("flex-1 flex gap-px")(

		// Button().
		// 	Class("rounded shadow bg-white").
		// 	Color(Blue).
		// 	Click(ctx.Call(collate.onSearch, reset).Replace(collate.Target)).
		// 	Render(Icon("fa fa-times")),

		Form("flex-1 flex bg-blue-800 rounded-l-lg", ctx.Send(collate.onSearch).AsSubmit(collate.Target))(
			IText("Search", query).
				Class("flex-1 p-1 w-72").
				ClassInput("cursor-pointer bg-white border-gray-300 hover:border-blue-500 block w-full p-3").
				Placeholder(ctx.Translate("Hľadať")).
				Render(""),

			Button().
				Submit().
				Class("rounded shadow bg-white").
				Color(Blue).
				Render(Icon("fa fa-fw fa-search")),
		),

		If(len(collate.Excel) > 0 || collate.OnExcel != nil, func() string {
			return Button().
				Color(Blue).
				Click(ctx.Call(collate.onXLS).None()).
				Render(Icon2("fa fa-download", "XLS"))
		}),

		If(len(collate.Filter) > 0, func() string {
			return Button().
				Submit().
				Class("rounded-r-lg shadow bg-white").
				Color(Blue).
				Click(fmt.Sprintf("window.document.getElementById('%s')?.classList.toggle('hidden');", collate.TargetFilter.Id)).
				Render(Icon3("fa fa-fw fa-chevron-down", "Filter"))
		}),
	)
}

func Sorting[T any](ctx *Context, collate *TCollate[T], query *TQuery) string {
	if len(collate.Sort) == 0 {
		return ""
	}

	return Div("flex gap-px")(
		Map(collate.Sort, func(sort *TField, index int) string {
			if sort.DB == "" {
				sort.DB = sort.Field
			}

			direction := ""
			color := GrayOutline
			field := strings.ToLower(sort.DB)
			query := strings.ToLower(query.Order)

			if strings.Contains(query, field) {
				if strings.Contains(query, "asc") {
					direction = "asc"
				} else {
					direction = "desc"
				}

				color = Purple
			}

			reverse := "desc"

			if direction == "desc" {
				reverse = "asc"
			}

			return Button().
				Class("rounded bg-white").
				Color(color).
				Click(ctx.Call(collate.onSort, TQuery{Order: sort.DB + " " + reverse}).Replace(collate.Target)).
				Render(
					Div("flex gap-2 items-center")(
						Iff(direction == "asc")(Icon("fa fa-fw fa-sort-amount-asc")),
						Iff(direction == "desc")(Icon("fa fa-fw fa-sort-amount-desc")),
						Iff(direction == "")(Icon("fa fa-fw fa-sort")),
						sort.Text,
					),
				)
		}),
	)
}

func Paging[T any](ctx *Context, collate *TCollate[T], result *TCollateResult[T]) string {
	if result.Filtered == 0 {
		return Empty(result)
	}

	size := len(result.Data)
	more := ctx.Translate("Load more items")
	count := ctx.Translate("Showing %d / %d of %d in total", size, result.Filtered, result.Total)

	if result.Filtered == result.Total {
		count = ctx.Translate("Showing %d / %d", size, result.Total)
	}

	return Div("flex items-center justify-center")(
		// showing information
		Div("mx-4 font-bold text-lg")(count),

		Div("flex gap-px flex-1 justify-end")(
			// reset
			If(collate.onReset != nil, func() string {
				return Button().
					Class("bg-white rounded-l").
					Color(PurpleOutline).
					Disabled(size == 0 || size <= int(collate.Limit)).
					Click(ctx.Call(collate.onReset).Replace(collate.Target)).
					Render(
						Icon("fa fa-fw fa-undo"),
						// Div("flex gap-2 items-center")(
						// 	Icon("fa fa-repeat"), reset,
						// ),
					)
			}),

			// load more
			If(collate.onResize != nil, func() string {
				return Button().
					Class("rounded-r").
					Color(Purple).
					Disabled(size >= int(result.Filtered)).
					Click(ctx.Call(collate.onResize).Replace(collate.Target)).
					Render(
						Div("flex gap-2 items-center")(
							Icon("fa fa-arrow-down"), more,
						),
					)
			}),
		),
	)
}
