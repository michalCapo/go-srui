import { Classes, INPUT, MD } from './ui.util';
import { Div } from './ui.shared';

export class ICheckbox {
  private data?: any; private name: string; private css = ''; private size = MD;
  private required = false;
  constructor(name: string, data?: any) { this.name = name; this.data = data; }
  Class(...v: string[]) { this.css = v.join(' '); return this; }
  Size(v: string) { this.size = v; return this; }
  Required(v = true) { this.required = v; return this; }
  Render(label: string): string {
    const checked = this.data ? Boolean((this.data as any)[this.name]) : false;
    return Div(this.css)(`<label class=\"flex items-center gap-2\"><input class=\"${Classes(INPUT, this.size)}\" type=\"checkbox\" name=\"${this.name}\" ${checked ? 'checked' : ''} ${this.required ? 'required' : ''}/> ${label}</label>`);
  }
}

export const Checkbox = (name: string, data?: any) => new ICheckbox(name, data);

