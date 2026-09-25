/**
 * 主题元数据公开端点
 * 原样返回 theme.json（$meta + schema），供 Server 定时任务/管端拉取（pull 模型）
 */
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

export default defineEventHandler(async event => {
  try {
    const themePath = resolve(process.cwd(), 'theme.json')
    const raw = JSON.parse(readFileSync(themePath, 'utf-8'))

    setHeader(event, 'Content-Type', 'application/json')
    setHeader(event, 'Cache-Control', 'no-store')
    return raw
  } catch {
    throw createError({ statusCode: 404, statusMessage: 'theme.json not found' })
  }
})
