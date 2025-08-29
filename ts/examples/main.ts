import * as ui from '../ui';
import { HelloContent } from './pages/hello';
import { ButtonContent } from './pages/button';
import { CounterContent } from './pages/counter';
import { LoginContent } from './pages/login';
import { ShowcaseContent } from './pages/showcase';

type Route = { Path: string; Title: string };
const routes: Route[] = [
  { Path: '/', Title: 'Hello' },
  { Path: '/button', Title: 'Button' },
  { Path: '/counter', Title: 'Counter' },
  { Path: '/login', Title: 'Login' },
  { Path: '/showcase', Title: 'Showcase' },
];

const app = ui.MakeApp('en');

const layout = (title: string, body: (ctx: ui.Context) => string): ui.Callable =>
  (ctx: ui.Context) => {
    const nav = ui.Div('bg-white shadow mb-6')(
      ui.Div('max-w-5xl mx-auto px-4 py-3 flex flex-wrap gap-2 items-center')(
        ui.Div('flex flex-wrap gap-2')(
          routes.map(r => ui.A('px-3 py-1 rounded hover:bg-gray-200', { href: r.Path }, ctx.Load(r.Path))(r.Title)).join(' '),
        ),
      ),
    );
    const content = body(ctx);
    return app.HTML(title, 'p-4 bg-gray-200 min-h-screen', nav + ui.Div('max-w-5xl mx-auto px-2')(content));
  };

app.Page('/', layout('Hello', HelloContent));
app.Page('/button', layout('Button', ButtonContent));
app.Page('/counter', layout('Counter', CounterContent));
app.Page('/login', layout('Login', LoginContent));
app.Page('/showcase', layout('Showcase', ShowcaseContent));

// Provide Node globals to satisfy TypeScript without @types/node
declare var require: any; declare var module: any;

// Start server when executed directly (node ts/examples/main.js)
if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  app.listen(1422);
}
