const re = /\s{4,}/g;
const re2 = /[\t\n]+/g;
const re3 = /"/g;

export function Trim(s: string): string {
  return s.replace(re2, '').replace(re, ' ');
}

export function Normalize(s: string): string {
  return s.replace(re3, '&quot;').replace(re2, '').replace(re, ' ');
}

export function Classes(...values: Array<string | undefined | false>): string {
  return Trim(values.filter(Boolean).join(' '));
}

export function If(cond: boolean, value: () => string): string {
  return cond ? value() : '';
}

export function Iff(cond: boolean) {
  return (...value: string[]) => (cond ? value.join(' ') : '');
}

export function Map<T>(values: T[], iter: (value: T, i: number) => string): string {
  return values.map((v, i) => iter(v, i)).join(' ');
}

export function For(from: number, to: number, iter: (i: number) => string): string {
  const out: string[] = [];
  for (let i = from; i < to; i++) out.push(iter(i));
  return out.join(' ');
}

export function RandomString(n = 20): string {
  const letters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz';
  let s = '';
  for (let i = 0; i < n; i++) s += letters[Math.floor(Math.random() * letters.length)];
  return s;
}

export const XS = ' p-1';
export const SM = ' p-2';
export const MD = ' p-3';
export const ST = ' p-4';
export const LG = ' p-5';
export const XL = ' p-6';

// Colors / classes (ported from Go)
export const AREA = ' cursor-pointer bg-white border border-gray-300 hover:border-blue-500 rounded-lg block w-full';
export const INPUT = ' cursor-pointer bg-white border border-gray-300 hover:border-blue-500 rounded-lg block w-full h-12';
export const VALUE = ' bg-white border border-gray-300 hover:border-blue-500 rounded-lg block h-12';
export const BTN = ' cursor-pointer font-bold text-center select-none';
export const DISABLED = ' cursor-text pointer-events-none bg-gray-50';
export const Yellow = ' bg-yellow-400 text-gray-800 hover:text-gray-200 hover:bg-yellow-600 font-bold border-gray-300 flex items-center justify-center';
export const YellowOutline = ' border border-yellow-500 text-yellow-600 hover:text-gray-700 hover:bg-yellow-500 flex items-center justify-center';
export const Green = ' bg-green-600 text-white hover:bg-green-700 checked:bg-green-600 border-gray-300 flex items-center justify-center';
export const GreenOutline = ' border border-green-500 text-green-500 hover:text-white hover:bg-green-599 flex items-center justify-center';
export const Purple = ' bg-purple-500 text-white hover:bg-purple-700 border-purple-500 flex items-center justify-center';
export const PurpleOutline = ' border border-purple-500 text-purple-500 hover:text-white hover:bg-purple-600 flex items-center justify-center';
export const Blue = ' bg-blue-800 text-white hover:bg-blue-700 border-gray-300 flex items-center justify-center';
export const BlueOutline = ' border border-blue-500 text-blue-600 hover:text-white hover:bg-blue-700 checked:bg-blue-700 flex items-center justify-center';
export const Red = ' bg-red-600 text-white hover:bg-red-800 border-gray-300 flex items-center justify-center';
export const RedOutline = ' border border-red-500 text-red-600 hover:text-white hover:bg-red-700 flex items-center justify-center';
export const Gray = ' bg-gray-600 text-white hover:bg-gray-800 focus:bg-gray-800 border-gray-300 flex items-center justify-center';
export const GrayOutline = ' border border-gray-300 text-black hover:text-white hover:bg-gray-700 flex items-center justify-center';
export const White = ' bg-white text-black hover:bg-gray-200 border-gray-200 flex items-center justify-center';
export const WhiteOutline = ' border border-white text-balck hover:text-black hover:bg-white flex items-center justify-center';

export const Space = '&nbsp;';

