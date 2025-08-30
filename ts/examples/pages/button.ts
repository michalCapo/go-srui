import * as ui from '../../ui.index';

export function ButtonContent(ctx: ui.Context): string {
  const buttonId = ui.Target();
  let show: ui.Callable;

  const button: ui.Callable = (ctx) => new ui.Button()
    .Click(ctx.Call(show).Replace(buttonId))
    .Class('rounded')
    .Color(ui.Blue)
    .Render('Click me');

  show = (ctx) => ui.Div('flex gap-2 items-center bg-red-500 rounded text-white p-px pl-4', buttonId)(
    'Clicked',
    new ui.Button()
      .Click(ctx.Call(button).Replace(buttonId))
      .Class('rounded')
      .Color(ui.Red)
      .Render('Hide me'),
  );

  return ui.Div('flex flex-row gap-4')(
    ui.Div('flex justify-start gap-4 items-center')(button(ctx)),
  );
}
