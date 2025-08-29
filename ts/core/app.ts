import { __load, __post, __stringify, __submit, Context, Target } from './context';
import { Classes, Trim } from './util';
import { Callable } from './types';

// Provide Node globals to satisfy TypeScript without @types/node
declare var require: any;

type Handler = Callable;

export class App {
  Language: string;
  HTMLHead: string[];
  HTMLBody: (cls: string) => string;

  private routes = new Map<string, Handler>();
  private stored = new Map<Handler, string>();

  constructor(defaultLanguage: string) {
    this.Language = defaultLanguage;
    this.HTMLHead = [
      `<meta charset="UTF-8">`,
      `<meta name="viewport" content="width=device-width, initial-scale=1.0">`,
      `<style>
        html { scroll-behavior: smooth; }
        .invalid, select:invalid, textarea:invalid, input:invalid {
          border-bottom-width: 2px; border-bottom-color: red; border-bottom-style: dashed;
        }
        @media (max-width: 768px) {
          input[type="date"] { max-width: 100% !important; width: 100% !important; min-width: 0 !important; box-sizing: border-box !important; overflow: hidden !important; }
          input[type="date"]::-webkit-datetime-edit { max-width: 100% !important; overflow: hidden !important; }
        }
      </style>`,
      `<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/tailwindcss/2.2.19/tailwind.min.css" integrity="sha512-wnea99uKIC3TJF7v4eKk4Y+lMz2Mklv18+r4na2Gn1abDRPPOeef95xTzdwGD9e6zXJBteMIhZ1+68QC5byJZw==" crossorigin="anonymous" referrerpolicy="no-referrer" />`,
      script(__stringify, __post, __submit, __load),
    ];
    this.HTMLBody = (cls: string) => {
      if (!cls) cls = 'bg-gray-200';
      const contentId = Target().id!;
      return `
        <!DOCTYPE html>
        <html lang="${this.Language}" class="${cls}">
          <head>__head__</head>
          <body id="${contentId}" class="relative">__body__</body>
        </html>`;
    };
  }

  HTML(title: string, bodyClass: string, ...body: string[]): string {
    const html = this.HTMLBody(Classes(bodyClass))
      .replace('__head__', [`<title>${title}</title>`, ...this.HTMLHead].join(''))
      .replace('__body__', body.join(''));
    return Trim(html);
  }

  private register(path: string, handler: Handler): void {
    if (!path) throw new Error('Path cannot be empty');
    if (this.routes.has(path)) throw new Error('Path already registered: ' + path);
    this.routes.set(path, handler);
    this.stored.set(handler, path);
  }

  Page(path: string, component: Handler): Handler {
    this.register(path, component);
    return component;
  }

  Action(uid: string, action: Handler): Handler {
    if (!uid.startsWith('/')) uid = '/' + uid;
    uid = uid.toLowerCase();
    if (Array.from(this.stored.values()).includes(uid)) {
      // return already registered handler if any
      for (const [fn, p] of this.stored.entries()) if (p === uid) return fn;
    }
    this.register(uid, action);
    return action;
  }

  Callable(action: Handler): Handler {
    if (this.stored.has(action)) return action;
    // create uid from function reference string
    const uid = '/' + (action.name || 'fn')
      .replace(/[.*()\[\]]/g, '')
      .replace(/[./:-]/g, '-')
      .toLowerCase();
    return this.Action(uid, action);
  }

  pathOf(handler: Handler): string | undefined {
    return this.stored.get(handler);
  }

  // Minimal HTTP server for Node (optional for running examples)
  listen(port = 1422) {
    // Lazy import to avoid hard Node requirement for consumers
    const http = require('http');
    const server = http.createServer(async (req: any, res: any) => {
      const url = req.url as string;
      const method = req.method as string;
      const path = url.split('?')[0];
      const handler = this.routes.get(path);
      if (!handler) { res.statusCode = 404; res.end('Not found'); return; }

      let body = '';
      req.on('data', (chunk: any) => (body += chunk));
      await new Promise(resolve => req.on('end', resolve));
      try { req.body = body ? JSON.parse(body) : undefined; } catch { req.body = undefined; }

      const ctx = new Context(this, req, res, 'sess-' + Math.random().toString(36).slice(2));
      res.setHeader('Content-Type', 'text/html; charset=utf-8');
      const html = handler(ctx);
      res.write(html);
      if (ctx.append.length) res.write(ctx.append.join(''));
      res.end();
    });
    server.listen(port, '0.0.0.0');
    // eslint-disable-next-line no-console
    console.log(`Listening on http://0.0.0.0:${port}`);
  }
}

export function MakeApp(defaultLanguage: string) { return new App(defaultLanguage); }

function script(...parts: string[]) { return `<script>${parts.join(' ')}</script>`; }
