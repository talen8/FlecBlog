<template>
  <Milkdown />
</template>

<script setup lang="ts">
import { watch, onUnmounted } from 'vue'
import { Milkdown, useEditor } from '@milkdown/vue'
import { Editor, rootCtx, defaultValueCtx, editorViewCtx, prosePluginsCtx, serializerCtx, parserCtx } from '@milkdown/kit/core'
import { commonmark } from '@milkdown/kit/preset/commonmark'
import { gfm } from '@milkdown/kit/preset/gfm'
import { history } from '@milkdown/kit/plugin/history'
import { clipboard } from '@milkdown/kit/plugin/clipboard'
import { indent } from '@milkdown/kit/plugin/indent'
import { Plugin, PluginKey } from '@milkdown/kit/prose/state'
import type { Ctx } from '@milkdown/kit/ctx'
import { content, setContent, charCount, isReadOnly } from '@/stores/editor'

const syncKey = new PluginKey('content-sync')

function createSyncPlugin(ctx: Ctx) {
  return new Plugin({
    key: syncKey,
    view: () => ({
      update: (view) => {
        try {
          const serializer = ctx.get(serializerCtx)
          const markdown = serializer(view.state.doc)
          if (markdown !== content.value) {
            setContent(markdown)
            const text = markdown.trim()
            charCount.value = markdown.replace(/\s/g, '').length
          }
        } catch {
          // serializer not ready yet
        }
      },
    }),
  })
}

const { get: getEditor } = useEditor((root) => {
  const editor = Editor.make()
  editor
    .config((ctx) => {
      ctx.set(rootCtx, root)
      ctx.set(defaultValueCtx, content.value)
    })
    .use(commonmark)
    .use(gfm)
    .use(history)
    .use(clipboard)
    .use(indent)
    .config((ctx) => {
      ctx.update(prosePluginsCtx, (plugins) => [...plugins, createSyncPlugin(ctx)])
    })
  return editor
})

// 外部内容变化时同步到编辑器
watch(content, (val) => {
  const editor = getEditor()
  if (!editor) return
  editor.action((ctx) => {
    const view = ctx.get(editorViewCtx)
    const serializer = ctx.get(serializerCtx)
    const current = serializer(view.state.doc)
    if (current === val) return
    const parser = ctx.get(parserCtx)
    const doc = parser(val)
    if (!doc) return
    const { state } = view
    const tr = (state as any).tr.replaceWith(0, state.doc.content.size, doc.content)
    view.dispatch(tr)
  })
})

// 只读状态变化
watch(isReadOnly, (val) => {
  const editor = getEditor()
  if (!editor) return
  editor.action((ctx) => {
    const view = ctx.get(editorViewCtx)
    view.dom.setAttribute('contenteditable', val ? 'false' : 'true')
  })
})

// 监听工具栏插入事件
function handleInsert(e: Event) {
  if (isReadOnly.value) return
  const editor = getEditor()
  if (!editor) return
  const { type } = (e as CustomEvent).detail

  const inserts: Record<string, { before: string; after: string; placeholder: string }> = {
    bold: { before: '**', after: '**', placeholder: '粗体' },
    italic: { before: '*', after: '*', placeholder: '斜体' },
    strikethrough: { before: '~~', after: '~~', placeholder: '删除线' },
    code: { before: '`', after: '`', placeholder: 'code' },
    heading: { before: '## ', after: '', placeholder: '标题' },
    quote: { before: '> ', after: '', placeholder: '引用' },
    ul: { before: '- ', after: '', placeholder: '列表项' },
    ol: { before: '1. ', after: '', placeholder: '列表项' },
    task: { before: '- [ ] ', after: '', placeholder: '任务' },
    link: { before: '[', after: '](url)', placeholder: '链接文本' },
    image: { before: '![', after: '](url)', placeholder: '图片描述' },
    codeblock: { before: '```\n', after: '\n```', placeholder: '代码' },
    table: { before: '| 列1 | 列2 | 列3 |\n| --- | --- | --- |\n', after: '', placeholder: '' },
    hr: { before: '\n---\n', after: '', placeholder: '' },
  }

  const fmt = inserts[type]
  if (!fmt) return

  editor.action((ctx) => {
    const view = ctx.get(editorViewCtx)
    const { from, to } = view.state.selection
    const selectedText = view.state.doc.textBetween(from, to, '\n')
    const selected = selectedText || fmt.placeholder
    const insert = fmt.before + selected + fmt.after
    const tr = (view.state as any).tr.insertText(insert, from, to)
    view.dispatch(tr)
    view.focus()
  })
}

window.addEventListener('editor:insert', handleInsert)
onUnmounted(() => {
  window.removeEventListener('editor:insert', handleInsert)
})
</script>

<style scoped lang="scss">
:deep([data-milkdown-root]) {
  height: 100%;
  overflow: auto;
}

:deep(.ProseMirror) {
  height: 100%;
  padding: 24px 32px;
  max-width: 800px;
  margin: 0 auto;
  outline: none;
  color: var(--text-primary);
  background: var(--bg-base);

  p {
    color: var(--text-regular);
    line-height: 1.8;
  }

  h1, h2, h3, h4, h5, h6 {
    color: var(--text-primary);
    border-bottom: 1px solid var(--border-secondary);
    padding-bottom: 8px;
  }

  ul, ol {
    color: var(--text-regular);
    padding-left: 20px;
  }

  li {
    color: var(--text-regular);
  }

  strong {
    color: var(--text-primary);
  }

  em {
    color: var(--text-secondary);
  }

  code {
    background: var(--bg-elevated);
    padding: 2px 6px;
    border-radius: 3px;
    color: var(--accent-red);
  }

  pre {
    background: var(--bg-elevated);
    padding: 16px;
    border-radius: 6px;
    overflow-x: auto;

    code {
      background: transparent;
      color: var(--text-primary);
    }
  }

  blockquote {
    border-left: 3px solid var(--accent-blue);
    padding-left: 12px;
    color: var(--text-secondary);
  }

  a {
    color: var(--accent-blue);
  }

  hr {
    border-color: var(--border-primary);
  }

  table {
    border-collapse: collapse;
    width: 100%;

    th, td {
      border: 1px solid var(--border-primary);
      padding: 8px 12px;
    }

    th {
      background: var(--bg-elevated);
      color: var(--text-primary);
    }

    td {
      color: var(--text-regular);
    }
  }
}
</style>
