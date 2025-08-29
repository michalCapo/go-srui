import { Attr } from '../core/types';
import { Classes, INPUT, MD } from '../core/util';
import { Div, Label } from './shared';

export interface OptionItem {
  id: string;
  value: string;
}

export class ISelect<T = any> {
  private data?: T;
  private name: string;
  private css = '';
  private cssLabel = '';
  private cssInput = '';
  private size = MD;
  private required = false;
  private disabled = false;
  private placeholder = '';
  private options: OptionItem[] = [];
  private target: Attr = {};

  constructor(name: string, data?: T) { this.name = name; this.data = data; }

  Class(...v: string[]) { this.css = v.join(' '); return this; }
  ClassLabel(...v: string[]) { this.cssLabel = v.join(' '); return this; }
  ClassInput(...v: string[]) { this.cssInput = v.join(' '); return this; }
  Size(v: string) { this.size = v; return this; }
  Required(v = true) { this.required = v; return this; }
  Disabled(v = true) { this.disabled = v; return this; }
  Options(values: OptionItem[]) { this.options = values; return this; }
  Placeholder(v: string) { this.placeholder = v; return this; }

  Render(label: string): string {
    const selected = this.data ? String((this.data as any)[this.name] ?? '') : '';
    const opts = [this.placeholder ? `<option value="">${this.placeholder}</option>` : '']
      .concat(this.options.map(o => `<option value="${o.id}" ${selected === o.id ? 'selected' : ''}>${o.value}</option>`))
      .join('');
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      `<select class="${Classes(INPUT, this.size, this.cssInput)}" id="${this.target.id ?? ''}" name="${this.name}" ${this.required ? 'required' : ''} ${this.disabled ? 'disabled' : ''}>${opts}</select>`
    );
  }
}

export function SelectInput<T>(name: string, data?: T) { return new ISelect<T>(name, data); }

