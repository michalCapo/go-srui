import { Attr } from '../core/types';
import { Classes, DISABLED, INPUT, MD, AREA } from '../core/util';
import { Div, Input, Label, Textarea } from './shared';

type ErrorLike = undefined | null;

class BaseInput<TData = any> {
  protected data?: TData;
  protected renderFn: (label: string) => string = () => '';
  protected placeholder = '';
  protected css = '';
  protected cssLabel = '';
  protected cssInput = '';
  protected autocomplete = '';
  protected size = MD;
  protected onclick = '';
  protected onchange = '';
  protected as = 'text';
  protected name: string;
  protected pattern = '';
  protected value = '';
  protected target: Attr = {};
  protected visible = true;
  protected required = false;
  protected disabled = false;
  protected readonly = false;

  constructor(name: string, data?: TData) {
    this.name = name;
    this.data = data;
  }

  Class(...v: string[]) { this.css = v.join(' '); return this; }
  ClassLabel(...v: string[]) { this.cssLabel = v.join(' '); return this; }
  ClassInput(...v: string[]) { this.cssInput = v.join(' '); return this; }
  Size(v: string) { this.size = v; return this; }
  Placeholder(v: string) { this.placeholder = v; return this; }
  Pattern(v: string) { this.pattern = v; return this; }
  Autocomplete(v: string) { this.autocomplete = v; return this; }
  Required(v = true) { this.required = v; return this; }
  Readonly(v = true) { this.readonly = v; return this; }
  Disabled(v = true) { this.disabled = v; return this; }
  Type(v: string) { this.as = v; return this; }
  Rows(v: number) { this.target.rows = v; return this; }
  Value(v: string) { this.value = v; return this; }
  Change(code: string) { this.onchange = code; return this; }
  Click(code: string) { this.onclick = code; return this; }
  If(v: boolean) { this.visible = v; return this; }

  protected resolveValue(): string {
    if (!this.data) return this.value;
    const val = (this.data as any)[this.name];
    if (val == null) return this.value;
    if (val instanceof Date) {
      if (this.as === 'date') return val.toISOString().slice(0, 10);
      if (this.as === 'time') return val.toISOString().slice(11, 16);
      if (this.as === 'datetime-local') return val.toISOString().slice(0, 16);
    }
    return String(val);
  }
}

export class IText extends BaseInput {
  constructor(name: string, data?: any) { super(name, data); this.as = 'text'; }
  Render(label: string, err?: ErrorLike): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onchange: this.onchange, onclick: this.onclick,
        required: this.required, disabled: this.disabled, value, pattern: this.pattern,
        placeholder: this.placeholder, autocomplete: this.autocomplete,
      }),
    );
  }
}

export class IPassword extends BaseInput {
  constructor(name: string, data?: any) { super(name, data); this.as = 'password'; }
  Render(label: string, err?: ErrorLike): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick,
        required: this.required, disabled: this.disabled, value, placeholder: this.placeholder,
      }),
    );
  }
}

export class IArea extends BaseInput {
  constructor(name: string, data?: any) { super(name, data); this.as = 'text'; }
  Render(label: string, err?: ErrorLike): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    const rows = this.target.rows ?? 5;
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Textarea(Classes(AREA, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick,
        required: this.required, disabled: this.disabled, readonly: this.readonly, placeholder: this.placeholder,
        rows,
      })(value),
    );
  }
}

export class INumber extends BaseInput {
  private min?: number; private max?: number; private step?: number; private valueFormat = '%v';
  constructor(name: string, data?: any) { super(name, data); this.as = 'number'; }
  Numbers(min?: number, max?: number, step?: number) { this.min = min; this.max = max; this.step = step; return this; }
  Format(fmt: string) { this.valueFormat = fmt; return this; }
  Render(label: string, err?: ErrorLike): string {
    if (!this.visible) return '';
    let value = this.resolveValue();
    if (this.valueFormat && value) {
      if (this.valueFormat.includes('%.2f')) {
        const n = Number(value); if (!Number.isNaN(n)) value = n.toFixed(2);
      }
    }
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick,
        required: this.required, disabled: this.disabled, value,
        min: this.min !== undefined ? String(this.min) : undefined,
        max: this.max !== undefined ? String(this.max) : undefined,
        step: this.step !== undefined ? String(this.step) : undefined,
        placeholder: this.placeholder,
      }),
    );
  }
}

export class IDate extends BaseInput {
  private min?: Date; private max?: Date;
  constructor(name: string, data?: any) { super(name, data); this.as = 'date'; }
  Dates(min?: Date, max?: Date) { this.min = min; this.max = max; return this; }
  Render(label: string): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    const min = this.min ? this.min.toISOString().slice(0, 10) : '';
    const max = this.max ? this.max.toISOString().slice(0, 10) : '';
    return Div(this.css + ' min-w-0')(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, 'min-w-0 max-w-full', this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick, onchange: this.onchange,
        required: this.required, disabled: this.disabled, value, min, max, placeholder: this.placeholder,
      }),
    );
  }
}

export class ITime extends BaseInput {
  private min?: Date; private max?: Date;
  constructor(name: string, data?: any) { super(name, data); this.as = 'time'; }
  Dates(min?: Date, max?: Date) { this.min = min; this.max = max; return this; }
  Render(label: string): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    const min = this.min ? this.min.toISOString().slice(11, 16) : '';
    const max = this.max ? this.max.toISOString().slice(11, 16) : '';
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick,
        required: this.required, disabled: this.disabled, value, min, max, placeholder: this.placeholder,
      }),
    );
  }
}

export class IDateTime extends BaseInput {
  private min?: Date; private max?: Date;
  constructor(name: string, data?: any) { super(name, data); this.as = 'datetime-local'; }
  Dates(min?: Date, max?: Date) { this.min = min; this.max = max; return this; }
  Render(label: string): string {
    if (!this.visible) return '';
    const value = this.resolveValue();
    const min = this.min ? this.min.toISOString().slice(0, 16) : '';
    const max = this.max ? this.max.toISOString().slice(0, 16) : '';
    return Div(this.css)(
      Label(this.cssLabel, { for: this.target.id, required: this.required })(label),
      Input(Classes(INPUT, this.size, this.cssInput, this.disabled && DISABLED), {
        id: this.target.id, name: this.name, type: this.as, onclick: this.onclick,
        required: this.required, disabled: this.disabled, value, min, max, placeholder: this.placeholder,
      }),
    );
  }
}

// Helper factory functions
export const Text = (name: string, data?: any) => new IText(name, data);
export const Password = (name: string, data?: any) => new IPassword(name, data);
export const Area = (name: string, data?: any) => new IArea(name, data);
export const NumberInput = (name: string, data?: any) => new INumber(name, data);
export const DateInput = (name: string, data?: any) => new IDate(name, data);
export const TimeInput = (name: string, data?: any) => new ITime(name, data);
export const DateTimeInput = (name: string, data?: any) => new IDateTime(name, data);

