<template>
  <div class="source-editor" ref="editorRoot"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { EditorView, keymap } from '@codemirror/view'
import { EditorState, Compartment } from '@codemirror/state'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import { markdown } from '@codemirror/lang-markdown'
import { oneDark } from '@codemirror/theme-one-dark'
import { content, setContent, charCount, isReadOnly } from '@/stores/editor'

const editorRoot = ref<HTMLElement>()
let editorView: EditorView | undefined
const readOnlyCompartment = new Compartment()

onMounted(() => {
  if (!editorRoot.value) return

  const state = EditorState.create({
    doc: content.value,
    extensions: [
      readOnlyCompartment.of(EditorState.readOnly.of(isReadOnly.value)),
      history(),
      keymap.of([...defaultKeymap, ...historyKeymap]),
      markdown(),
      EditorView.lineWrapping,
      oneDark,
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          const doc = update.state.doc.toString()
          setContent(doc)
          charCount.value = doc.replace(/\s/g, '').length
        }
      }),
    ],
  })

  editorView = new EditorView({
    state,
    parent: editorRoot.value,
  })

  // 初始化统计
  charCount.value = content.value.replace(/\s/g, '').length

  window.addEventListener('editor:insert', handleInsert)
})

onUnmounted(() => {
  editorView?.destroy()
  window.removeEventListener('editor:insert', handleInsert)
})

// 外部内容变化时同步
watch(content, (val) => {
  if (!editorView) return
  const current = editorView.state.doc.toString()
  if (current !== val) {
    editorView.dispatch({
      changes: { from: 0, to: current.length, insert: val },
    })
  }
})

// 只读状态变化
watch(isReadOnly, (val) => {
  if (!editorView) return
  editorView.dispatch({
    effects: readOnlyCompartment.reconfigure(EditorState.readOnly.of(val)),
  })
})

function handleInsert(e: Event) {
  if (isReadOnly.value || !editorView) return
  const { type } = (e as CustomEvent).detail
  const { from, to } = editorView.state.selection.main
  const doc = editorView.state.doc.toString()
  const selected = doc.slice(from, to)

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

  const text = selected || fmt.placeholder
  const insert = fmt.before + text + fmt.after
  editorView.dispatch({
    changes: { from, to, insert },
    selection: { anchor: from + fmt.before.length, head: from + fmt.before.length + text.length },
  })
  editorView.focus()
}
</script>

<style scoped lang="scss">
.source-editor {
  height: 100%;
  overflow: hidden;

  :deep(.cm-editor) {
    height: 100%;
    background: var(--bg-base);

    .cm-scroller {
      padding: 24px 32px;
      font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', Consolas, monospace;
      font-size: 14px;
      line-height: 1.8;
    }

    .cm-content {
      max-width: 800px;
      margin: 0 auto;
    }

    &.cm-focused {
      outline: none;
    }
  }
}
</style>
