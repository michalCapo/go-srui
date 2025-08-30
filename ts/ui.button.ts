import { Attr } from './ui.types';
import { BTN, Classes, DISABLED, MD } from './ui.util';
import { A, Div } from './ui.elements';

export class Button {
  private size = MD;
  private color = '';
  private onclick = '';
  private css = '';
  private as: 'div' | 'button' | 'a' = 'div';
  private target: Attr = {};
  private visible = true;
  private disabled = false;
  private extra: Attr[] = [];

  constructor(...attr: Attr[]) {
    this.extra = attr;
  }

  Submit() { this.as = 'button'; this.extra.push({ type: 'submit' }); return this; }
  Reset() { this.as = 'button'; this.extra.push({ type: 'reset' }); return this; }
  If(v: boolean) { this.visible = v; return this; }
  Disabled(v: boolean) { this.disabled = v; return this; }
  Class(...v: string[]) { this.css = v.join(' '); return this; }
  Color(v: string) { this.color = v; return this; }
  Size(v: string) { this.size = v; return this; }
  Click(code: string) { this.onclick = code; return this; }
  Href(v: string) { this.as = 'a'; this.extra.push({ href: v }); return this; }

  Render(text: string): string {
    if (!this.visible) return '';
    const cls = Classes(BTN, this.size, this.color, this.css, this.disabled && DISABLED + ' opacity-25');

    if (this.as === 'a') {
      return A(cls, ...this.extra, { id: this.target.id })(text);
    }

    if (this.as === 'div') {
      return Div(cls, ...this.extra, { id: this.target.id, onclick: this.onclick })(text);
    }

    // button element
    return `<button ${[...this.extra, { id: this.target.id, onclick: this.onclick, class: cls }]
      .map(a => Object.entries(a).map(([k, v]) => v !== undefined && v !== false ? `${k}=\"${v === true ? k : v}\"` : '').filter(Boolean).join(' ')).join(' ')}>${text}</button>`;
  }
}

export function ButtonEl(...attr: Attr[]) { return new Button(...attr); }

