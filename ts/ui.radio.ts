import { Attr } from './ui.types';
import { Classes, INPUT, MD } from './ui.util';
import { Div, Label } from './ui.shared';

export interface AOption { id: string; value: string; }

export class IRadio {
  private data?: any;
  private name: string;
  private css = '';
  private cssLabel = '';
  private size = MD;
  private valueSet = '';
  private target: Attr = {};

  constructor(name: string, data?: any) { this.name = name; this.data = data; }
  Class(...v: string[]) { this.css = v.join(' '); return this; }
  ClassLabel(...v: string[]) { this.cssLabel = v.join(' '); return this; }
  Size(v: string) { this.size = v; return this; }
  Value(v: string) { this.valueSet = v; return this; }

  Render(label: string): string {
    const selected = this.data ? String((this.data as any)[this.name] ?? '') : '';
    const checked = selected === this.valueSet ? 'checked' : '';
    return Div(this.css)(
      `<label class=\"${this.cssLabel}\">\n         <input class=\"${Classes(INPUT, this.size)}\" type=\"radio\" name=\"${this.name}\" value=\"${this.valueSet}\" ${checked}/> ${label}\n       </label>`
    );
  }
}

export class IRadioButtons {
  private data?: any; private name: string; private css = '';
  private options: AOption[] = [];
  constructor(name: string, data?: any) { this.name = name; this.data = data; }
  Options(v: AOption[]) { this.options = v; return this; }
  Class(...v: string[]) { this.css = v.join(' '); return this; }
  Render(label: string): string {
    const selected = this.data ? String((this.data as any)[this.name] ?? '') : '';
    const items = this.options.map(o => `<label class=\"px-3 py-2 border rounded ${selected === o.id ? 'bg-blue-700 text-white' : ''}\"><input type=\"radio\" name=\"${this.name}\" value=\"${o.id}\" ${selected === o.id ? 'checked' : ''}/> ${o.value}</label>`).join(' ');
    return Div(this.css)(Label('font-bold')(`${label}`), Div('flex gap-2 flex-wrap')(items));
  }
}

export const Radio = (name: string, data?: any) => new IRadio(name, data);
export const RadioButtons = (name: string, data?: any) => new IRadioButtons(name, data);

