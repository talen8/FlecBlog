/**
 * 主题元数据同步插件
 * 启动时读取 theme.json，调用 /themes/_sync 注册到 Server
 */
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

const MAX_RETRIES = 5
const RETRY_DELAY_MS = 2000

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

export default defineNitroPlugin(async () => {
  const config = useRuntimeConfig()
  const apiUrl = config.public.apiUrl as string

  if (!apiUrl) {
    console.warn('[theme-sync] apiUrl 未配置，跳过同步')
    return
  }

  // 读取 theme.json
  const themePath = resolve(process.cwd(), 'theme.json')
  const fullSchema = JSON.parse(readFileSync(themePath, 'utf-8'))

  // 提取 $meta，剩余作为 schema
  const { $meta, ...schema } = fullSchema

  if (!$meta?.slug) {
    console.warn('[theme-sync] theme.json 中缺少 $meta.slug，跳过同步')
    return
  }

  // 重试机制：应对 Server 尚未就绪的情况
  for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
    try {
      await $fetch('/themes/_sync', {
        method: 'POST',
        baseURL: apiUrl,
        body: {
          slug: $meta.slug,
          name: $meta.name,
          version: $meta.version,
          author: $meta.author,
          description: $meta.description,
          license: $meta.license,
          repo: $meta.repo,
          schema,
        },
      })

      console.log(`[theme-sync] 主题 ${$meta.slug} 同步成功`)
      return
    } catch (error) {
      console.error(`[theme-sync] 同步失败 (${attempt}/${MAX_RETRIES}):`, error)
      if (attempt < MAX_RETRIES) {
        await sleep(RETRY_DELAY_MS * attempt)
      }
    }
  }

  console.error(`[theme-sync] 主题 ${$meta.slug} 同步失败，已重试 ${MAX_RETRIES} 次`)
})
