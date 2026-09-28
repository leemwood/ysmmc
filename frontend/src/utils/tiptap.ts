// TipTap JSON → 安全 HTML 轻量渲染器。
// 只输出白名单标签与自构造属性，文本节点一律 HTML 转义，图片仅放行 https。

const ESC: Record<string, string> = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  "'": '&#39;',
};

function esc(s: string): string {
  return s.replace(/[&<>"']/g, (c) => ESC[c]);
}

function safeSrc(src: string): string | null {
  if (typeof src !== 'string' || !/^https:\/\//i.test(src)) return null;
  return src;
}

type Node = {
  type?: string;
  text?: string;
  attrs?: Record<string, unknown>;
  content?: Node[];
  marks?: { type?: string }[];
};

function renderInline(nodes: Node[] | undefined): string {
  if (!nodes) return '';
  let out = '';
  for (const n of nodes) {
    if (n.type === 'text' && typeof n.text === 'string') {
      let piece = esc(n.text);
      for (const m of n.marks || []) {
        if (m.type === 'bold') piece = `<strong>${piece}</strong>`;
        else if (m.type === 'italic') piece = `<em>${piece}</em>`;
        else if (m.type === 'code') piece = `<code>${piece}</code>`;
      }
      out += piece;
    } else if (n.type === 'hardBreak') {
      out += '<br>';
    } else if (n.type === 'image') {
      const src = safeSrc(String(n.attrs?.src ?? ''));
      if (src) {
        out += `<img src="${esc(src)}" alt="${esc(String(n.attrs?.alt ?? ''))}" loading="lazy">`;
      }
    } else if (n.content) {
      out += renderInline(n.content);
    }
  }
  return out;
}

function renderBlock(node: Node): string {
  switch (node.type) {
    case 'paragraph':
      return `<p>${renderInline(node.content)}</p>`;
    case 'heading': {
      const level = Math.min(Math.max(Number(node.attrs?.level) || 2, 2), 4);
      return `<h${level}>${renderInline(node.content)}</h${level}>`;
    }
    case 'bulletList':
      return `<ul>${(node.content || []).map((li) => `<li>${renderInline(li.content)}</li>`).join('')}</ul>`;
    case 'orderedList':
      return `<ol>${(node.content || []).map((li) => `<li>${renderInline(li.content)}</li>`).join('')}</ol>`;
    case 'blockquote':
      return `<blockquote>${renderBlockChildren(node.content)}</blockquote>`;
    case 'codeBlock':
      return `<pre><code>${renderInline(node.content)}</code></pre>`;
    case 'horizontalRule':
      return '<hr>';
    default:
      return node.content ? renderBlockChildren([node]) : renderInline([node]);
  }
}

function renderBlockChildren(nodes: Node[] | undefined): string {
  if (!nodes) return '';
  return nodes.map(renderBlock).join('');
}

/**
 * 把 TipTap JSON 字符串/对象渲染为受限 HTML。
 * 解析失败或输入为纯文本时，返回转义后的段落。
 */
export function tiptapToHtml(content: unknown): string {
  let doc: Node | null = null;
  if (typeof content === 'string' && content.trim().startsWith('{')) {
    try {
      doc = JSON.parse(content) as Node;
    } catch {
      doc = null;
    }
  } else if (content && typeof content === 'object') {
    doc = content as Node;
  }
  if (doc && doc.type === 'doc' && Array.isArray(doc.content) && doc.content.length) {
    return renderBlockChildren(doc.content);
  }
  if (typeof content === 'string' && content.trim()) {
    return content
      .split(/\n{2,}|\r\n\r\n/)
      .map((p) => `<p>${esc(p).replace(/\n/g, '<br>')}</p>`)
      .join('');
  }
  return '';
}
