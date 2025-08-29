import * as ui from '../../ui';

export class TCounter { constructor(public Count: number) {} }

export function Counter(count: number) { return new TCounter(count); }

export function CounterContent(ctx: ui.Context): string {
  return ui.Div('flex flex-row gap-4')(renderCounter(ctx, Counter(3)));
}

function renderCounter(ctx: ui.Context, counter: TCounter): string {
  const target = ui.Target();

  const decrement: ui.Callable = (ctx) => {
    ctx.Body(counter);
    counter.Count--; if (counter.Count < 0) counter.Count = 0;
    return renderCounter(ctx, counter);
  };

  const increment: ui.Callable = (ctx) => {
    ctx.Body(counter);
    counter.Count++;
    return renderCounter(ctx, counter);
  };

  return ui.Div('flex gap-2 items-center bg-purple-500 rounded text-white p-px', target)(
    new ui.Button().Click(ctx.Call(decrement, counter).Replace(target)).Class('rounded-l px-5').Render('-'),
    ui.Div('text-2xl')(`${counter.Count}`),
    new ui.Button().Click(ctx.Call(increment, counter).Replace(target)).Class('rounded-r px-5').Render('+'),
  );
}

