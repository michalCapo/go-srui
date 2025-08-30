import { Attr } from './ui.types';
import { Classes } from './ui.util';

function attributes(...attrs: Attr[]): string {
  const result: string[] = [];
  for (const a of attrs) {
    if (!a) continue;
    if (a.id) result.push(`id="${a.id}"`);
    if (a.href) result.push(`href="${a.href}"`);
    if (a.alt) result.push(`alt="${a.alt}"`);
    if (a.title) result.push(`title="${a.title}"`);
    if (a.src) result.push(`src="${a.src}"`);
    if (a.for) result.push(`for="${a.for}"`);
    if (a.type) result.push(`type="${a.type}"`);
    if (a.class) result.push(`class="${a.class}"`);
    if (a.style) result.push(`style="${a.style}"`);
    if (a.onclick) result.push(`onclick="${a.onclick}"`);
    if (a.onchange) result.push(`onchange="${a.onchange}"`);
    if (a.onsubmit) result.push(`onsubmit="${a.onsubmit}"`);
    if (a.value !== undefined) result.push(`value="${a.value ?? ''}"`);
    if (a.checked) result.push(`checked="${a.checked}"`);
    if (a.selected) result.push(`selected="${a.selected}"`);
    if (a.name) result.push(`name="${a.name}"`);
    if (a.placeholder) result.push(`placeholder="${a.placeholder}"`);
    if (a.autocomplete) result.push(`autocomplete="${a.autocomplete}"`);
    if (a.pattern) result.push(`pattern="${a.pattern}"`);
    if (a.cols) result.push(`cols="${a.cols}"`);
    if (a.rows) result.push(`rows="${a.rows}"`);
    if (a.width) result.push(`width="${a.width}"`);
    if (a.height) result.push(`height="${a.height}"`);
    if (a.min) result.push(`min="${a.min}"`);
    if (a.max) result.push(`max="${a.max}"`);
    if (a.target) result.push(`target="${a.target}"`);
    if (a.step) result.push(`step="${a.step}"`);
    if (a.required) result.push('required="required"');
    if (a.disabled) result.push('disabled="disabled"');
    if (a.readonly) result.push('readonly="readonly"');
  }
  return result.join(' ');
}

function open(tag: string) {
  return (css: string, ...attr: Attr[]) =>
    (...elements: string[]) => {
      const final = [...attr, { class: Classes(css) }];
      return `<${tag} ${attributes(...final)}>${elements.join(' ')}</${tag}>`;
    };
}

function closed(tag: string) {
  return (css: string, ...attr: Attr[]) => {
    const final = [...attr, { class: Classes(css) }];
    return `<${tag} ${attributes(...final)}/>`;
  };
}

export const I = open('i');
export const A = open('a');
export const P = open('p');
export const Div = open('div');
export const Span = open('span');
export const Form = open('form');
export const Textarea = open('textarea');
export const Select = open('select');
export const Option = open('option');
export const List = open('ul');
export const ListItem = open('li');
export const Canvas = open('canvas');
export const Img = closed('img');
export const Input = closed('input');

export { attributes };

