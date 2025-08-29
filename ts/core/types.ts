export type Swap = 'inline' | 'outline' | 'none';
export type ActionType = 'POST' | 'FORM';

export interface Attr {
  onclick?: string;
  onchange?: string;
  onsubmit?: string;
  step?: string;
  id?: string;
  href?: string;
  title?: string;
  alt?: string;
  type?: string;
  class?: string;
  style?: string;
  name?: string;
  value?: string;
  checked?: string;
  for?: string;
  src?: string;
  selected?: string;
  pattern?: string;
  placeholder?: string;
  autocomplete?: string;
  max?: string;
  min?: string;
  target?: string;
  rows?: number;
  cols?: number;
  width?: number;
  height?: number;
  disabled?: boolean;
  required?: boolean;
  readonly?: boolean;
}

export interface BodyItem {
  name: string;
  type: string;
  value: string;
}

export type Callable = (ctx: import('./context').Context) => string;

export interface CSSMut {
  orig: string;
  set?: string;
  append: string[];
}

