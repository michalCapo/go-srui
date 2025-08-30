import { Attr } from './ui.types';
import { Div as DivCore, Input as InputCore, Textarea as TextareaCore } from './ui.elements';

// Re-export commonly used tags for UI components
export const Div = DivCore;
export const Input = InputCore;
export const Textarea = TextareaCore;

// Simple label wrapper
export function Label(css: string, ...attr: Attr[]) {
  return (text: string) => `<label class=\"${css}\" ${attrs(attr)}>${text}</label>`;
}

function attrs(attr: Attr[]): string {
  return attr
    .map(a => Object.entries(a)
      .map(([k, v]) => v !== undefined && v !== false ? `${k}=\"${v === true ? k : v}\"` : '')
      .filter(Boolean).join(' '))
    .join(' ');
}

