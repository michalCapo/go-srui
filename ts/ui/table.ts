export class SimpleTable {
  private cols: number; private css: string; private rows: string[][] = [];
  private headerClass: Record<number, string> = {};
  constructor(cols: number, css = '') { this.cols = cols; this.css = css; }
  Class(col: number, css: string) { this.headerClass[col] = css; return this; }
  Field(value: string) {
    if (this.rows.length === 0 || this.rows[this.rows.length - 1].length === this.cols) this.rows.push([]);
    this.rows[this.rows.length - 1].push(value);
    return this;
  }
  Render(): string {
    const rowsHtml = this.rows.map((row, i) => {
      const tag = i === 0 ? 'th' : 'td';
      return `<tr>${row.map((cell, j) => {
        const cls = i === 0 ? (this.headerClass[j] || '') : '';
        return `<${tag} class="${cls} p-2 border-b">${cell}</${tag}>`;}).join('')}</tr>`;
    }).join('');
    return `<table class="${this.css}"><tbody>${rowsHtml}</tbody></table>`;
  }
}

export function Table(cols: number, css?: string) { return new SimpleTable(cols, css); }

