# 多阶段构建 Dockerfile for Nuxt 3

# 阶段 1: 构建阶段
FROM node:22-alpine AS builder

WORKDIR /app

# 复制 package 文件
COPY package*.json ./

# 安装依赖
RUN npm ci

# 复制源代码
COPY . .

# 构建应用
RUN npm run build

# 阶段 2: 运行阶段
FROM node:22-alpine

WORKDIR /app

# 复制构建产物和必要文件
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/package*.json ./

# 设置环境变量
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=3000

# 暴露端口
EXPOSE 3000

# 启动应用
CMD ["node", ".output/server/index.mjs"]
