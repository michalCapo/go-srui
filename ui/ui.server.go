package ui

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/websocket"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Method = func(*Context) string

// type Partial struct {
// 	Method *Method
// 	Path   string
// }

var (
	stored = make(map[*Method]string)
	mu     sync.Mutex
)

type BodyItem struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CSS struct {
	Orig   string
	Set    string
	Append []string
}

func (c *CSS) Value() string {
	if len(c.Set) > 0 {
		return c.Set
	}

	c.Append = append(c.Append, c.Orig)
	return Classes(c.Append...)
}

type Swap string

const (
	OUTLINE Swap = "outline"
	INLINE  Swap = "inline"
	NONE    Swap = "none"
)

type ActionType string

const (
	POST ActionType = "POST"
	FORM ActionType = "FORM"
)

type Context struct {
	App       *App
	Request   *http.Request
	Response  http.ResponseWriter
	SessionId string
	append    []string
}

type TSession struct {
	DB        *gorm.DB `gorm:"-"`
	SessionId string
	Name      string
	Data      datatypes.JSON
}

func (TSession) TableName() string {
	return "_session"
}

func (session *TSession) Load(data any) {
	temp := &TSession{}

	err := session.DB.Where("session_id = ? and name = ?", session.SessionId, session.Name).Take(temp).Error
	if err != nil {
		return
	}

	err = json.Unmarshal(temp.Data, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (session *TSession) Save(output any) {
	data, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}

	temp := &TSession{
		SessionId: session.SessionId,
		Name:      session.Name,
		Data:      data,
	}

	session.DB.Where("session_id = ? and name = ?", session.SessionId, session.Name).Save(temp)
}

func (ctx *Context) Ip() string {
	return ctx.Request.RemoteAddr
}

func (ctx *Context) Session(db *gorm.DB, name string) *TSession {
	return &TSession{
		DB:        db,
		Name:      name,
		SessionId: ctx.SessionId,
	}
}

// func (ctx *Context) SessionSave(db *gorm.DB, name string, output any) {
// 	data, err := json.Marshal(output)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	session := &TSession{
// 		SessionId: ctx.SessionId,
// 		Name:      name,
// 		Data:      data,
// 	}

// 	db.Where("session_id = ? and name = ?", session.SessionId, session.Name).Save(session)
// }

func (ctx *Context) Body(output any) error {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}

	var data []BodyItem
	if len(body) > 0 {
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
	}

	// structValue := reflect.ValueOf(output).Elem()

	for _, item := range data {
		structFieldValue, err := PathValue(output, item.Name)
		// structFieldValue := structValue.FieldByName(item.Name)
		if err != nil {
			fmt.Println("Error getting field", item.Name, err)
			continue
		}

		// if !structFieldValue.IsValid() { continue }

		if !structFieldValue.CanSet() {
			continue
		}

		val := reflect.ValueOf(item.Value)

		if structFieldValue.Type() != val.Type() {
			switch item.Type {
			case "date":
				t, err := time.Parse("2006-01-02", item.Value)
				if err != nil {
					fmt.Println("Error parsing date", err)
					continue
				}
				if structFieldValue.Type() == reflect.TypeOf(gorm.DeletedAt{}) {
					val = reflect.ValueOf(gorm.DeletedAt{Time: t, Valid: true})
				} else {
					val = reflect.ValueOf(t)
				}

			case "bool", "checkbox":
				val = reflect.ValueOf(item.Value == "true")

			case "radio", "string":
				val = reflect.ValueOf(item.Value)

			case "time":
				t, err := time.Parse("15:04", item.Value)
				if err != nil {
					fmt.Println("Error parsing time", err)
					continue
				}
				val = reflect.ValueOf(t)

			case "Time":
				t, err := time.Parse("2006-01-02 15:04:05 -0700 UTC", item.Value)
				if err != nil {
					fmt.Println("Error parsing time", err)
				}
				val = reflect.ValueOf(t)

			case "uint":
				cleanedValue := strings.ReplaceAll(item.Value, "_", "")
				n, err := strconv.ParseUint(cleanedValue, 10, 64)
				if err != nil {
					fmt.Println("Error parsing number", err)
					continue
				}
				val = reflect.ValueOf(uint(n))

			case "int":
				cleanedValue := strings.ReplaceAll(item.Value, "_", "")
				n, err := strconv.ParseInt(cleanedValue, 10, 64)
				if err != nil {
					fmt.Println("Error parsing number", err)
					continue
				}
				val = reflect.ValueOf(int(n))

			case "int64":
				cleanedValue := strings.ReplaceAll(item.Value, "_", "")
				n, err := strconv.ParseInt(cleanedValue, 10, 64)
				if err != nil {
					fmt.Println("Error parsing number", err)
					continue
				}
				val = reflect.ValueOf(int64(n))

			case "number":
				cleanedValue := strings.ReplaceAll(item.Value, "_", "")
				n, err := strconv.Atoi(cleanedValue)
				if err != nil {
					fmt.Println("Error parsing number", err)
					continue
				}
				val = reflect.ValueOf(n)

			case "float64":
				cleanedValue := strings.ReplaceAll(item.Value, "_", "")
				f, err := strconv.ParseFloat(cleanedValue, 64)
				if err != nil {
					fmt.Println("Error parsing float64", err)
					continue
				}
				val = reflect.ValueOf(f)

			case "datetime-local":
				t, err := time.Parse("2006-01-02T15:04", item.Value)
				if err != nil {
					fmt.Println("Error parsing datetime-local", err)
					continue
				}
				val = reflect.ValueOf(t)

			// case "text":
			// 	val = reflect.ValueOf(item.Value)

			case "":
				continue

			case "Model": // gorm.Model
				continue

			default:
				fmt.Println("Skipping (name;type;value):", item.Name, ";", item.Type, ";", item.Value)
				continue
			}
		}

		// fmt.Println("Setting", item.Name, "to", item.Value)
		structFieldValue.Set(val)
	}

	return nil
}

func (ctx *Context) Action(uid string, action Method) **Method {
	if ctx.App == nil {
		panic("App is nil, cannot register component. Did you set the App field in Context?")
	}
	return ctx.App.Action(uid, action)
}

func (ctx *Context) Callable(action Method) **Method {
	if ctx.App == nil {
		panic("App is nil, cannot create callable. Did you set the App field in Context?")
	}

	return ctx.App.Callable(action)
}

func (ctx *Context) Post(as ActionType, swap Swap, action *Action) string {
	path, ok := stored[action.Method]

	if !ok {
		a := reflect.ValueOf(*action.Method).String()
		// funcName := runtime.FuncForPC(reflect.ValueOf(action.Method).Pointer()).Name()
		panic(fmt.Sprintf("Function '%s' probably not registered. Cannot make call to this function.", a))
	}

	var body []BodyItem

	for _, item := range action.Values {
		v := reflect.ValueOf(item)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := range v.NumField() {
			field := v.Field(i)
			fieldName := v.Type().Field(i).Name
			fieldType := field.Type().Name()
			fieldValue := field.Interface()

			body = append(body, BodyItem{
				Name:  fieldName,
				Type:  fieldType,
				Value: fmt.Sprintf("%v", fieldValue),
			})
		}
	}

	values := "[]"

	if len(body) > 0 {
		temp, err := json.Marshal(body)

		if err == nil {
			values = string(temp)
		}
	}

	if as == FORM {
		return Normalize(fmt.Sprintf(`__submit(event, "%s", "%s", "%s", %s) `, swap, action.Target.Id, path, values))
	}

	return Normalize(fmt.Sprintf(`__post(event, "%s", "%s", "%s", %s) `, swap, action.Target.Id, path, values))
}

type Actions struct {
	Render   func(target Attr) string
	Replace  func(target Attr) string
	None     func() string
	AsSubmit func(target Attr, swap ...Swap) Attr
	AsClick  func(target Attr, swap ...Swap) Attr
}

type Submits struct {
	Render  func(target Attr) Attr
	Replace func(target Attr) Attr
}

// func (ctx *Context) ClickTo(method ComponentMethod, values ...any) AttrActions {
// 	return AttrActions{
// 		Replace: func(target Attr) Attr {
// 			target.Swap = "outline"
// 			return Attr{OnClick: ctx.Post(Action{Method: method, Target: target, Values: values})}
// 		},
// 		Render: func(target Attr) Attr {
// 			target.Swap = "inline"
// 			return Attr{OnClick: ctx.Post(Action{Method: method, Target: target, Values: values})}
// 		},
// 	}
// }
// func (ctx *Context) ChangeTo(method ComponentMethod, values ...any) AttrActions {
// 	return AttrActions{
// 		Replace: func(target Attr) Attr {
// 			target.Swap = "outline"
// 			return Attr{OnChange: ctx.Post(Action{Method: method, Target: target, Values: values})}
// 		},
// 		Render: func(target Attr) Attr {
// 			target.Swap = "inline"
// 			return Attr{OnChange: ctx.Post(Action{Method: method, Target: target, Values: values})}
// 		},
// 	}
// }

func swapize(swap ...Swap) Swap {
	if len(swap) > 0 {
		return swap[0]
	}

	return INLINE
}

func (ctx *Context) Submit(method **Method, values ...any) Submits {
	return Submits{
		Render: func(target Attr) Attr {
			return Attr{OnSubmit: ctx.Post(FORM, INLINE, &Action{Method: *method, Target: target, Values: values})}
		},
		Replace: func(target Attr) Attr {
			return Attr{OnSubmit: ctx.Post(FORM, OUTLINE, &Action{Method: *method, Target: target, Values: values})}
		},
	}
}

func (ctx *Context) Send(method **Method, values ...any) Actions {
	return Actions{
		// Replace: func(target Attr) string {
		// 	target.Swap = "outline"
		// 	return ctx.Post(Action{Type: FORM, Method: method, Target: target, Values: values})
		// },
		Render: func(target Attr) string {
			return ctx.Post(FORM, INLINE, &Action{Method: *method, Target: target, Values: values})
		},
		Replace: func(target Attr) string {
			return ctx.Post(FORM, OUTLINE, &Action{Method: *method, Target: target, Values: values})
		},
		None: func() string {
			return ctx.Post(FORM, NONE, &Action{Method: *method, Values: values})
		},
		AsSubmit: func(target Attr, swap ...Swap) Attr {
			return Attr{OnSubmit: ctx.Post(FORM, swapize(swap...), &Action{Method: *method, Target: target, Values: values})}
		},
		AsClick: func(target Attr, swap ...Swap) Attr {
			return Attr{OnClick: ctx.Post(FORM, swapize(swap...), &Action{Method: *method, Target: target, Values: values})}
		},
	}
}

func (ctx *Context) Call(method **Method, values ...any) Actions {
	return Actions{
		// Body: func() (*Context, string) {
		// 	return ctx, ctx.Post&Action{Method: method, Values: values})
		// },
		// Replace: func(target Attr) string {
		// 	target.Swap = "outline"
		// 	return ctx.Post&Action{Method: method, Target: target, Values: values})
		// },
		Render: func(target Attr) string {
			return ctx.Post(POST, INLINE, &Action{Method: *method, Target: target, Values: values})
		},
		Replace: func(target Attr) string {
			return ctx.Post(POST, OUTLINE, &Action{Method: *method, Target: target, Values: values})
		},
		None: func() string {
			return ctx.Post(POST, NONE, &Action{Method: *method, Values: values})
		},
		AsSubmit: func(target Attr, swap ...Swap) Attr {
			return Attr{OnSubmit: ctx.Post(POST, swapize(swap...), &Action{Method: *method, Target: target, Values: values})}
		},
		AsClick: func(target Attr, swap ...Swap) Attr {
			return Attr{OnClick: ctx.Post(POST, swapize(swap...), &Action{Method: *method, Target: target, Values: values})}
		},
	}
}

func (ctx *Context) Load(href string) Attr {
	return Attr{OnClick: Normalize(fmt.Sprintf(`__load("%s")`, href))}
}

func (ctx *Context) Reload() string {
	// return Normalize("<html><!DOCTYPE html><body><script>window.location.reload();</script></body></html>")
	return Normalize("<script>window.location.reload();</script>")
}

func (ctx *Context) Redirect(href string) string {
	// return Normalize(fmt.Sprintf("<html><!DOCTYPE html><body><script>window.location.href = '%s';</script></body></html>", href))
	return Normalize(fmt.Sprintf("<script>window.location.href = '%s';</script>", href))
}

func displayMessage(ctx *Context, message string, color string) {
	ctx.append = append(ctx.append,
		Trim((`<script>
            (function() {
                const el = document.getElementById("__messages__");
                if(el == null) {
                    const loader = document.createElement("div");
                    loader.id = "__messages__";
                    loader.classList = "fixed top-0 right-0 p-2 z-40";
                    document.body.appendChild(loader);
                }
            })();
        </script>`)),

		Trim(fmt.Sprintf(`<script>
            (function () {
                const el = document.getElementById("__messages__");
                if(el != null) {
                    const loader = document.createElement("div");
                    loader.classList = "p-4 m-2 rounded text-center border border-gray-700 shadow-xl text-xl text-center w-64 %s";
                    loader.innerHTML = "%s";
                    el.appendChild(loader);
					setTimeout(() => el.removeChild(loader), 5000);
                }
            })();
        </script>`, color, Normalize(message))),
	)
}

func (ctx *Context) Success(message string) {
	displayMessage(ctx, message, "bg-green-700 text-white")
}

func (ctx *Context) Error(message string) {
	displayMessage(ctx, message, "bg-red-700 text-white")
}

func (ctx *Context) DownloadAs(file *io.Reader, content_type string, name string) error {
	// Read the file content into a byte slice
	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		log.Println(err)
		return err
	}

	// Encode the byte slice to a base64 string
	fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)

	ctx.append = append(ctx.append,
		Trim(fmt.Sprintf(`<script>
            (function () {
                const byteCharacters = atob("%s");
                const byteNumbers = new Array(byteCharacters.length);
                for (let i = 0; i < byteCharacters.length; i++) {
                    byteNumbers[i] = byteCharacters.charCodeAt(i);
                }
                const byteArray = new Uint8Array(byteNumbers);
                const blob = new Blob([byteArray], { type: "%s" });
                const url = URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.href = url;
                a.download = "%s";
                a.click();
                URL.revokeObjectURL(url);
            })();
        </script>`, fileBase64, content_type, name)),
	)

	return nil
}

func (ctx *Context) Translate(message string, val ...any) string {
	return fmt.Sprintf(message, val...)
}

// type Component Method

func RandomString(n ...int) string {
	if len(n) == 0 {
		return RandomString(20)
	}

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n[0])
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func cacheControlMiddleware(next http.Handler, maxAge time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Cache-Control header
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(maxAge.Seconds())))

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

type App struct {
	Lanugage string
	HtmlBody func(string) string
	HtmlHead []string
}

// func (app *App) Path(method Method) string {
// 	funcName := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
// 	md5 := md5.Sum([]byte(funcName))
// 	path := hex.EncodeToString(md5[:])
// 	return "/" + strings.ToLower(path)
// }

func (app *App) Register(httpMethod string, path string, method *Method) string {
	if path == "" || method == nil {
		panic("Path and Method cannot be empty")
	}

	funcName := runtime.FuncForPC(reflect.ValueOf(*method).Pointer()).Name()

	if funcName == "" {
		panic("Method cannot be empty")
	}

	_, ok := stored[method]
	if ok {
		panic(fmt.Sprintf("Path %s is already registered", funcName))
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	for _, value := range stored {
		if value == path {
			panic(fmt.Sprintf("Path already exists in registry: %s -> %s", path, funcName))
		}
	}

	mu.Lock()
	stored[method] = path
	mu.Unlock()

	fmt.Println("Registering: ", httpMethod, path, " -> ", funcName)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(httpMethod, r.Method) {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var sessionId string
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sessionId = RandomString(30)
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionId,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				// Expires:  time.Now().Add(time.Hour * 24 * 30),
			})
		} else {
			sessionId = cookie.Value
		}

		ctx := &Context{
			App:       app,
			Request:   r,
			Response:  w,
			SessionId: sessionId,
			append:    []string{},
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte((*method)(ctx)))

		if len(ctx.append) > 0 {
			w.Write([]byte(strings.Join(ctx.append, "")))
		}
	})

	return path
}

func (app *App) Page(path string, component Method) **Method {
	found, ok := pool[path]
	if ok {
		return &found
	}

	found = &component
	app.Register("GET POST", path, found)
	pool[path] = found

	return &found
}

var pool = make(map[string]*Method)

func (app *App) Action(uid string, action Method) **Method {
	found, ok := pool[uid]
	if ok {
		return &found
	}

	found = &action
	app.Register("POST", uid, found)
	pool[uid] = found

	return &found
}

func (app *App) Callable(action Method) **Method {
	uid := runtime.FuncForPC(reflect.ValueOf(action).Pointer()).Name()
	uid = strings.ToLower(uid)
	uid = strings.ReplaceAll(uid, ".", "_")
	uid = strings.ReplaceAll(uid, "-", "_")
	uid = strings.ReplaceAll(uid, "/", "_")
	uid = strings.ReplaceAll(uid, ":", "_")
	uid = strings.ReplaceAll(uid, "*", "")
	uid = strings.ReplaceAll(uid, "(", "")
	uid = strings.ReplaceAll(uid, ")", "")

	// uid = fmt.Sprintf("%x", md5.Sum([]byte(uid)))

	found, ok := pool[uid]
	if ok {
		return &found
	}

	found = &action
	app.Register("POST", uid, found)
	pool[uid] = found

	return &found
}

func (app *App) Assets(assets embed.FS, path string, maxAge time.Duration) {
	path = strings.TrimPrefix(path, "/")
	http.Handle("/"+path, cacheControlMiddleware(http.FileServer(http.FS(assets)), maxAge))
}

func (app *App) Favicon(assets embed.FS, path string, maxAge time.Duration) {
	path = strings.TrimPrefix(path, "/")
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		file, err := assets.ReadFile(path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(maxAge.Seconds())))
		w.Write(file)
	})
}

func (app *App) Listen(port string) {
	log.Println("Listening on http://0.0.0.0" + port)

	if err := http.ListenAndServe(port, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error:", err)
	}
}

func (app *App) Autoreload() {
	app.HtmlHead = append(app.HtmlHead, `
		<script>
			const socket = new WebSocket('ws://' + window.location.host + '/live');
			socket.addEventListener('close', function (event) {
				document.body.innerHTML += '<div class="fixed inset-0 z-40 opacity-75 bg-gray-800"></div>';
				document.body.innerHTML += '<div class="fixed z-50 top-6 left-6 p-6 text-white bg-red-700 rounded border border-gray-500 uppercase font-bold">Offline</div>';
				setInterval(() => {
					fetch('/').then(() => window.location.reload())
				}, 1000);
			});
		</script>
	`)

	http.Handle("/live", websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		for {
			time.Sleep(10 * time.Second)
			ws.Write([]byte("ok"))
		}
	}))

	// http.HandleFunc("/live2", func(w http.ResponseWriter, r *http.Request) { for { time.Sleep(time.Minute) } })
}

func (app *App) Description(description string) {
	app.HtmlHead = append(app.HtmlHead, `<meta name="description" content="`+description+`">`)
}

// func (app *App) Script(script string) {
// 	app.HtmlHead = append(app.HtmlHead, `<script>`+script+`</script>`)
// }

func (app *App) Html(title string, class string, body ...string) string {
	head := []string{
		`<title>` + title + `</title>`,
	}

	head = append(head, app.HtmlHead...)

	html := app.HtmlBody(class)
	html = strings.ReplaceAll(html, "__lang__", app.Lanugage)
	html = strings.ReplaceAll(html, "__head__", strings.Join(head, " "))
	html = strings.ReplaceAll(html, "__body__", strings.Join(body, " "))

	return Trim(html)
}

var __post = Trim(` 
    function __post(event, swap, target_id, path, values) {
		const el = event.target;
		const name = el.getAttribute("name");
		const type = el.getAttribute("type");
		const value = el.value;

		let body = values; 
		if (name != null) {
			body = body.filter(element => element.name !== name);
			body.push({ name, type, value });
		}

		let loader;
		let loading = setTimeout(() => {
			loader = document.createElement("div");
			loader.classList = "fixed inset-0 flex gap-4 items-center justify-center z-50 bg-white opacity-75 font-bold text-3xl";
			loader.innerHTML = "Loading ...";
			document.body.appendChild(loader);
		}, 100);

		fetch(path, {method: "POST", body: JSON.stringify(body)})
			.then(html => html.text())
			.then(function (html) {
				const parser = new DOMParser();
				const doc = parser.parseFromString(html, 'text/html');
				const scripts = [...doc.body.querySelectorAll('script'), ...doc.head.querySelectorAll('script')];

				for (let i = 0; i < scripts.length; i++) {
					const newScript = document.createElement('script');
					newScript.textContent = scripts[i].textContent;
					document.body.appendChild(newScript);
				}

				const el = document.getElementById(target_id);
				if (el != null) {
					if (swap === "inline") {
						el.innerHTML = html;
					} else if(swap === "outline") {
						el.outerHTML = html;
					}
				}
			})
			.finally(function() {
				clearTimeout(loading);
				if(loader) {
					document.body.removeChild(loader);
				}
			});
    }
`)

var __stringify = Trim(`
    function __stringify(values) {
        const result = {};

        values.forEach(item => {
            const nameParts = item.name.split('.');
            let currentObj = result;
        
            for (let i = 0; i < nameParts.length - 1; i++) {
                const part = nameParts[i];
                if (!currentObj[part]) {
                    currentObj[part] = {};
                }
                currentObj = currentObj[part];
            }
        
            const lastPart = nameParts[nameParts.length - 1];

            switch(item.type) {
                case 'date':
                case 'time':
                case 'Time':
                case 'datetime-local':
                    currentObj[lastPart] = new Date(item.value);    
                    break;
                case 'float64':
                    currentObj[lastPart] = parseFloat(item.value);
                    break;
                case 'bool':
                case 'checkbox':
                    currentObj[lastPart] = item.value === 'true';
                    break;
                default:
                    currentObj[lastPart] = item.value;
            }
        });

        return JSON.stringify(result);
    }
`)

var __submit = Trim(`
    function __submit(event, swap, target_id, path, values) {
        event.preventDefault(); 

        const el = event.target;
        const tag = el.tagName.toLowerCase();
        const form = tag === "form" ? el : el.closest("form");
        const id = form.getAttribute("id");
        let body = values; 

        let found = Array.from(document.querySelectorAll('[form=' + id + '][name]'));

        if (found.length === 0) {
            found = Array.from(form.querySelectorAll('[name]'));
        };

        found.forEach((item) => {
            const name = item.getAttribute("name");
            const type = item.getAttribute("type");
            let value = item.value;
            
            if (type === 'checkbox') {
                value = String(item.checked)
            }

            if(name != null) {
                body = body.filter(element => element.name !== name);
                body.push({ name, type, value });
            }
        });

        let loader;
        let loading = setTimeout(() => {
            loader = document.createElement("div");
            loader.classList = "fixed inset-0 flex gap-4 items-center justify-center z-50 bg-white opacity-75 font-bold text-3xl";
            loader.innerHTML = "Loading ...";
            document.body.appendChild(loader);
        }, 100);

        fetch(path, {method: "POST", body: JSON.stringify(body)})
            .then(html => html.text())
			.then(function (html) {
				const parser = new DOMParser();
				const doc = parser.parseFromString(html, 'text/html');
				const scripts = [...doc.body.querySelectorAll('script'), ...doc.head.querySelectorAll('script')];

				for (let i = 0; i < scripts.length; i++) {
					const newScript = document.createElement('script');
					newScript.textContent = scripts[i].textContent;
					document.body.appendChild(newScript);
				}

				const el = document.getElementById(target_id);
				if (el != null) {
					if (swap === "inline") {
						el.innerHTML = html;
					} else if(swap === "outline") {
						el.outerHTML = html;
					}
				}
			})
            .finally(function() {
                clearTimeout(loading);
                if(loader) {
                    document.body.removeChild(loader);
                }
            });
    }
`)

var __load = Trim(`
    function __load(href) {
		event.preventDefault(); 

		let loader;
		let loading = setTimeout(() => {
			loader = document.createElement("div");
			loader.classList = "fixed inset-0 flex gap-4 items-center justify-center z-50 bg-white opacity-75 font-bold text-3xl";
			loader.innerHTML = "Loading ...";
			document.body.appendChild(loader);
		}, 100);

		fetch(href, {method: "GET"})
			.then(html => html.text())
			.then(function (html) {
				const parser = new DOMParser();
				const doc = parser.parseFromString(html, 'text/html');

				document.title = doc.title;
				document.body.innerHTML = doc.body.innerHTML;

				const scripts = [...doc.body.querySelectorAll('script'), ...doc.head.querySelectorAll('script')];
				for (let i = 0; i < scripts.length; i++) {
					const newScript = document.createElement('script');
					newScript.textContent = scripts[i].textContent;
					document.body.appendChild(newScript);
				}

				window.history.pushState({}, doc.title, href);
			})
			.finally(function() {
				clearTimeout(loading);
				if(loader) {
					document.body.removeChild(loader);
				}
			});
    }
`)

var ContentId = Target()

// var HtmlClass = "bg-gray-200 h-full"
// body_class = "max-w-4xl mx-auto p-4 lg:p-8 gap-8 h-full"

func MakeApp(default_language string) *App {

	return &App{
		Lanugage: default_language,
		HtmlHead: []string{
			`<meta charset="UTF-8">`,
			`<meta name="viewport" content="width=device-width, initial-scale=1.0">`,
			`<style>
				html {
					scroll-behavior: smooth;
				}
				.invalid, select:invalid, textarea:invalid, input:invalid {
					border-bottom-width: 2px;
					border-bottom-color: red;
					border-bottom-style: dashed;
				}
			</style>`,
			`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.2.19/tailwind.min.css" integrity="sha512-wnea99uKIC3TJF7v4eKk4Y+lMz2Mklv18+r4na2Gn1abDRPPOeef95xTzdwGD9e6zXJBteMIhZ1+68QC5byJZw==" crossorigin="anonymous" referrerpolicy="no-referrer" />`,
			Script(__stringify, __post, __submit, __load),
		},
		HtmlBody: func(class string) string {
			if class == "" {
				// class = "bg-gray-200 p-4 lg:p-8 gap-8 h-full"
				class = "bg-gray-200 h-full"
			}

			return fmt.Sprintf(`
				<!DOCTYPE html>
				<html lang="__lang__" class="%s">
					<head>__head__</head>
					<body id="%s" class="relative">__body__</body>
				</html>
			`, class, ContentId.Id)
		},
	}
}
