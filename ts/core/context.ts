import { ActionType, Attr, BodyItem, Callable, Swap } from './types';
import { Normalize, RandomString, Trim } from './util';

export class Context {
  app: import('./app').App;
  req: any;
  res: any;
  sessionID: string;
  append: string[] = [];

  constructor(app: import('./app').App, req: any, res: any, sessionID: string) {
    this.app = app;
    this.req = req;
    this.res = res;
    this.sessionID = sessionID;
  }

  Body<T extends object>(output: T): void {
    try {
      const data = this.req?.body as BodyItem[] | undefined;
      if (!Array.isArray(data)) return;
      for (const item of data) {
        setPath(output as any, item.name, coerce(item.type, item.value));
      }
    } catch {
      /* noop */
    }
  }

  Callable(method: Callable): Callable {
    return this.app.Callable(method);
  }

  Action(uid: string, action: Callable): Callable {
    return this.app.Action(uid, action);
  }

  Post(as: ActionType, swap: Swap, action: { method: Callable; target?: Attr; values?: any[] }): string {
    const path = this.app.pathOf(action.method);
    if (!path) throw new Error('Function not registered.');

    const body: BodyItem[] = [];
    for (const item of action.values || []) {
      if (item == null) continue;
      const entries = Object.entries(item);
      for (const [name, value] of entries) {
        body.push({ name, type: typeOf(value), value: valueToString(value) });
      }
    }

    const values = body.length > 0 ? JSON.stringify(body) : '[]';

    if (as === 'FORM') {
      return Normalize(`__submit(event, "${swap}", "${action.target?.id ?? ''}", "${path}", ${values}) `);
    }
    return Normalize(`__post(event, "${swap}", "${action.target?.id ?? ''}", "${path}", ${values}) `);
  }

  Send(method: Callable, ...values: any[]) {
    const callable = this.Callable(method);
    return {
      Render: (target: Attr) => this.Post('FORM', 'inline', { method: callable, target, values }),
      Replace: (target: Attr) => this.Post('FORM', 'outline', { method: callable, target, values }),
      None: () => this.Post('FORM', 'none', { method: callable, values }),
    };
  }

  Call(method: Callable, ...values: any[]) {
    const callable = this.Callable(method);
    return {
      Render: (target: Attr) => this.Post('POST', 'inline', { method: callable, target, values }),
      Replace: (target: Attr) => this.Post('POST', 'outline', { method: callable, target, values }),
      None: () => this.Post('POST', 'none', { method: callable, values }),
    };
  }

  Submit(method: Callable, ...values: any[]) {
    const callable = this.Callable(method);
    return {
      Render: (target: Attr): Attr => ({ onsubmit: this.Post('FORM', 'inline', { method: callable, target, values }) }),
      Replace: (target: Attr): Attr => ({ onsubmit: this.Post('FORM', 'outline', { method: callable, target, values }) }),
      None: (): Attr => ({ onsubmit: this.Post('FORM', 'none', { method: callable, values }) }),
    };
  }

  Load(href: string): Attr { return { onclick: Normalize(`__load("${href}")`) }; }
  Reload(): string { return Normalize('<script>window.location.reload();</script>'); }
  Redirect(href: string): string { return Normalize(`<script>window.location.href = '${href}';</script>`); }

  Success(message: string) { displayMessage(this, message, 'bg-green-700 text-white'); }
  Error(message: string) { displayMessage(this, message, 'bg-red-700 text-white'); }
}

function displayMessage(ctx: Context, message: string, color: string) {
  ctx.append.push(
    Trim(`<script>(function(){const el=document.getElementById("__messages__");if(el==null){const n=document.createElement("div");n.id="__messages__";n.classList="fixed top-0 right-0 p-2 z-40";document.body.appendChild(n);}})();</script>`),
  );
  ctx.append.push(
    Trim(`<script>(function(){const el=document.getElementById("__messages__");if(el!=null){const n=document.createElement("div");n.classList="p-4 m-2 rounded text-center border border-gray-700 shadow-xl text-xl text-center w-64 ${color}";n.innerHTML="${Normalize(message)}";el.appendChild(n);setTimeout(()=>el.removeChild(n),5000);}})();</script>`),
  );
}

export function Target(): Attr { return { id: 'i' + RandomString(15) }; }

// Serialize helpers
function typeOf(v: any): string {
  if (v instanceof Date) return 'Time';
  const t = typeof v;
  if (t === 'number') return Number.isInteger(v) ? 'int' : 'float64';
  if (t === 'boolean') return 'bool';
  if (t === 'string') return 'string';
  return 'string';
}

function valueToString(v: any): string {
  if (v instanceof Date) return v.toUTCString();
  return String(v);
}

function coerce(type: string, value: string): any {
  switch (type) {
    case 'date':
    case 'time':
    case 'Time':
    case 'datetime-local':
      return new Date(value);
    case 'float64':
      return parseFloat(value);
    case 'bool':
    case 'checkbox':
      return value === 'true';
    case 'int':
    case 'int64':
    case 'number':
      return parseInt(value, 10);
    default:
      return value;
  }
}

function setPath(obj: any, path: string, value: any) {
  const parts = path.split('.');
  let current = obj;
  for (let i = 0; i < parts.length - 1; i++) {
    const part = parts[i];
    if (!(part in current) || typeof current[part] !== 'object') current[part] = {};
    current = current[part];
  }
  current[parts[parts.length - 1]] = value;
}

// Client-side helpers injected to pages
export const __post = Trim(`
    function __post(event, swap, target_id, path, body) {
        event.preventDefault(); 
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
                    if (swap === "inline") { el.innerHTML = html; }
                    else if (swap === "outline") { el.outerHTML = html; }
                }
            })
            .finally(function() {
                clearTimeout(loading);
                if(loader) { document.body.removeChild(loader); }
            });
    }
`);

export const __stringify = Trim(`
    function __stringify(values) {
        const result = {};
        values.forEach(item => {
            const nameParts = item.name.split('.');
            let currentObj = result;
            for (let i = 0; i < nameParts.length - 1; i++) {
                const part = nameParts[i];
                if (!currentObj[part]) { currentObj[part] = {}; }
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
`);

export const __submit = Trim(`
    function __submit(event, swap, target_id, path, values) {
        event.preventDefault(); 
        const el = event.target;
        const tag = el.tagName.toLowerCase();
        const form = tag === "form" ? el : el.closest("form");
        const id = form.getAttribute("id");
        let body = values; 
        let found = Array.from(document.querySelectorAll('[form=' + id + '][name]'));
        if (found.length === 0) { found = Array.from(form.querySelectorAll('[name]')); };
        found.forEach((item) => {
            const name = item.getAttribute("name");
            const type = item.getAttribute("type");
            let value = item.value;
            if (type === 'checkbox') { value = String(item.checked) }
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
                    if (swap === "inline") { el.innerHTML = html; }
                    else if(swap === "outline") { el.outerHTML = html; }
                }
            })
            .finally(function() {
                clearTimeout(loading);
                if(loader) { document.body.removeChild(loader); }
            });
    }
`);

export const __load = Trim(`
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
                if(loader) { document.body.removeChild(loader); }
            });
    }
`);
