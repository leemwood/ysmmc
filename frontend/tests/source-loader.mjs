import { readFileSync, existsSync, statSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import ts from 'typescript'
import { parse, compileScript, compileTemplate } from '@vue/compiler-sfc'

const src = new URL('../src/', import.meta.url)

export async function resolve(specifier, context, next) {
  if (specifier.startsWith('@/')) {
    const base = new URL(specifier.slice(2), src)
    for (const suffix of ['', '.ts', '/index.ts']) {
      const candidate = new URL(base.href + suffix)
      if (existsSync(candidate) && statSync(candidate).isFile()) {
        return next(candidate.href, context)
      }
    }
  }
  return next(specifier, context)
}

export async function load(url, context, next) {
  if (!url.endsWith('.ts') && !url.endsWith('.vue')) return next(url, context)
  let source = readFileSync(fileURLToPath(url), 'utf8')
  if (url.endsWith('.vue')) {
    const { descriptor } = parse(source, { filename: fileURLToPath(url) })
    // 内存渲染器没有浏览器的 HTML 解析能力。
    const compilerOptions = { hoistStatic: false }
    if (descriptor.script || descriptor.scriptSetup) {
      source = compileScript(descriptor, {
        id: url,
        inlineTemplate: true,
        templateOptions: { compilerOptions },
      }).content
    } else {
      source = compileTemplate({ source: descriptor.template.content, filename: url, id: url, compilerOptions }).code
        + '\nexport default { render }'
    }
  }
  const { outputText } = ts.transpileModule(source, {
    compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.ESNext },
  })
  return { format: 'module', shortCircuit: true, source: outputText }
}
