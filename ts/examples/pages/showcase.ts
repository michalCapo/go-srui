import * as ui from '../../ui.index';

class DemoForm {
  Name = '';
  Email = '';
  Phone = '';
  Password = '';
  Age = 0;
  Price = 0;
  Bio = '';
  Gender = '';
  Country = '';
  Agree = false;
  BirthDate = new Date();
  AlarmTime = new Date();
  Meeting = new Date();
}

const demoTarget = ui.Target();

export function ShowcaseContent(ctx: ui.Context): string {
  const form = new DemoForm();
  return ui.Div('max-w-full sm:max-w-6xl mx-auto flex flex-col gap-6 w-full')(
    ui.Div('text-3xl font-bold')('TS-SRUI Component Showcase'),
    render(ctx, form, undefined),
  );
}

function render(ctx: ui.Context, f: DemoForm, err?: Error): string {
  const Submit: ui.Callable = (ctx) => { ctx.Body(f); ctx.Success('Form submitted successfully'); return render(ctx, f, undefined); };

  const countries = [ '', 'USA', 'Slovakia', 'Germany', 'Japan' ].map(x => ({ id: x, value: x || 'Select...' }));
  const genders = [ { id: 'male', value: 'Male' }, { id: 'female', value: 'Female' }, { id: 'other', value: 'Other' } ];

  return ui.Div('grid gap-4 sm:gap-6 lg:grid-cols-2 items-start w-full', demoTarget)(
    ui.Form('flex flex-col gap-4 bg-white p-6 rounded-lg shadow w-full', demoTarget, ctx.Submit(Submit).Replace(demoTarget))(
      ui.Div('text-xl font-bold')('Component Showcase Form'),
      err ? ui.Div('text-red-600 p-4 rounded text-center border-4 border-red-600 bg-white')(err.message) : '',

      new ui.IText('Name', f).Required().Render('Name'),
      new ui.IText('Email', f).Required().Render('Email'),
      new ui.IText('Phone', f).Render('Phone'),
      new ui.IPassword('Password').Required().Render('Password'),

      new ui.INumber('Age', f).Numbers(0, 120, 1).Render('Age'),
      new ui.INumber('Price', f).Format('%.2f').Render('Price (USD)'),
      new ui.IArea('Bio', f).Rows(4).Render('Short Bio'),

      ui.Div('block sm:hidden')(
        ui.Div('text-sm font-bold')('Gender'),
        new ui.IRadio('Gender', f).Value('male').Render('Male'),
        new ui.IRadio('Gender', f).Value('female').Render('Female'),
        new ui.IRadio('Gender', f).Value('other').Render('Other'),
      ),
      ui.Div('hidden sm:block overflow-x-auto')(
        new ui.IRadioButtons('Gender', f).Options(genders).Render('Gender'),
      ),
      new ui.ISelect('Country', f).Options(countries).Placeholder('Select...').Render('Country'),
      new ui.ICheckbox('Agree', f).Required().Render('I agree to the terms'),

      new ui.IDate('BirthDate', f).Render('Birth Date'),
      new ui.ITime('AlarmTime', f).Render('Alarm Time'),
      new ui.IDateTime('Meeting', f).Render('Meeting (Local)'),

      ui.Div('flex gap-2 mt-2')(
        new ui.Button().Submit().Color(ui.Blue).Class('rounded').Render('Submit'),
        new ui.Button().Reset().Color(ui.Gray).Class('rounded').Render('Reset'),
      ),
    ),

    ui.Div('flex flex-col gap-4 w-full')(
      ui.Div('bg-white p-6 rounded-lg shadow flex flex-col gap-2 w-full')(
        ui.Div('text-xl font-bold')('Buttons & Colors'),
        ui.Div('grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2')(
          new ui.Button().Color(ui.Blue).Class('rounded w-full').Render('Blue'),
          new ui.Button().Color(ui.Green).Class('rounded w-full').Render('Green'),
          new ui.Button().Color(ui.Red).Class('rounded w-full').Render('Red'),
          new ui.Button().Color(ui.Purple).Class('rounded w-full').Render('Purple'),
          new ui.Button().Color(ui.Yellow).Class('rounded w-full').Render('Yellow'),
          new ui.Button().Color(ui.Gray).Class('rounded w-full').Render('Gray'),
        ),
      ),

      ui.Div('bg-white p-6 rounded-lg shadow flex flex-col gap-3 w-full')(
        ui.Div('text-xl font-bold')('Counter (Actions)'),
        // simple inline counter demo
        (() => {
          const c = { Count: 2 } as any;
          const target = ui.Target();
          const render = (): string => ui.Div('flex gap-2 items-center bg-purple-500 rounded text-white p-px', target)(
            new ui.Button().Click(ctx.Call(() => { c.Count--; if (c.Count < 0) c.Count = 0; return render(); }).Replace(target)).Class('rounded-l px-5').Render('-'),
            ui.Div('text-2xl')(`${c.Count}`),
            new ui.Button().Click(ctx.Call(() => { c.Count++; return render(); }).Replace(target)).Class('rounded-r px-5').Render('+'),
          );
          return render();
        })(),
      ),
    ),
  );
}
