<div align="center">
  <img src=".github/images/logo.png" alt="FlecBlog Logo" width="92" height="92" />

  <h1>FlecBlog</h1>

  <p>
    A modern full-stack blogging system with taste, clarity, and a strong sense of identity.
  </p>

  <p>
    Instead of piling on noisy features, FlecBlog focuses on doing the essentials well:
    content, management, and presentation, all polished into one cohesive experience.
  </p>

  <p>
    <a href="https://blog.talen.top">Live Demo</a> /
    <a href="https://ccnlf8xcz6k3.feishu.cn/wiki/space/7618178485001046989">Documentation</a> /
    <a href="https://github.com/talen8/FlecBlog/issues/new">Report an Issue</a> /
    <a href="https://qm.qq.com/q/Zzm9XN6lOi">Community</a>
  </p>

  <p>
    <a href="./README.md">简体中文</a> / English
  </p>

  <p>
    <img src="https://img.shields.io/badge/Monorepo-FlecBlog-111111?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAACXBIWXMAAAsTAAALEwEAmpwYAAACDUlEQVR4nO3XS4iNYRgH8BNRNhYmkctMkpXJSslth1hiKxsbRmqoQU2NjY2J1GwsbGgsSG5FQ0kozMrKQhZqUCbkMsjdT2/znvo6fd93vjnznSPlX2dx3vN/nv//vT3PeyqV//hXgZt4jSEcRjeuYASf42ckjnVHzlCMuV2GgVGNY7QMA5cmYeBiGQb2T8JAT6Oii7E97ufdmqSrsCZFLG1870REZ+MQntSZVRBZW3D8VBHhqTiIj8rHiXriM3FL87A7T3w67jVR/BPm5Bno11yczRPvwNcmG/iG9r81+yr60sRn4I3W4E6agc1ah1dpBo7lBIRq1o4DIbgEA9/TDAznBKxO8NpiW50I3uECBvEstRvieU6C33iAdYlaUdsLshDeDW0JnWnYmWZgrECyn+hKXNkPBWL6I38DtmJK1hX8UnBGwcTKGNNXgB9W70Xi+74sA48j4WGBpMOJnhGeWBPB0SwDVyOhq855qGJLjFuAXXhbIOYl5mUZ6I2kc1has2xpeIpFifjOOmU85OtMFQ/AEvyKScIjZC6uFdjfR9XE4aGRwTuPWZV6ME4MGEyMbcSNKJaFI5HbU3NYL2N9XeEqsDBxtXprfpuPHTiJ+7GgjMXzsjxy9mAA28LZqDQC4zP+EU2cLrR0ZQObEoXpfewTKzKLSJNMdMS/Vcm9v94yAwkjy3A8PsvPlC3wB5cqDes3vUmxAAAAAElFTkSuQmCC&style=flat-square" alt="Monorepo" />
    <img src="https://img.shields.io/badge/Server-Go%201.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
    <img src="https://img.shields.io/badge/Admin-Vue%203-42B883?style=flat-square&logo=vuedotjs&logoColor=white" alt="Vue 3" />
    <img src="https://img.shields.io/badge/Blog-Nuxt%204.2.2-00DC82?style=flat-square&logo=nuxt&logoColor=white" alt="Nuxt 4" />
  </p>
</div>

<p align="center">
  <img src=".github/images/preview.png" alt="FlecBlog Preview" />
</p>

## Overview

FlecBlog is a three-part blogging system built around a complete content workflow.

By separating the public blog, the admin console, and the backend service into clearly defined layers, it stays expressive on the surface while remaining maintainable under the hood.

| Module | Stack | Purpose |
| --- | --- | --- |
| `server` | Go 1.25 / Gin / GORM / PostgreSQL | Backend services, authentication, APIs, data management, and scheduled jobs |
| `admin` | Vue 3 / Element Plus / Vite | Content management, dashboards, the editor, and admin operations |
| `blog` | Nuxt 4.2.2 / Vue 3 / SCSS | Public-facing blog, SSR, SEO, and reading experience |

**Why FlecBlog**

- Built as a complete engineering project, not just a polished demo page
- Clean separation between the admin panel and the public blog
- Covers the essentials you expect from a modern blog system: posts, comments, links, and analytics
- A strong fit for personal brands, content-driven websites, and long-term product expansion

## Preview

| Home Page | Article Page |
| --- | --- |
| ![Blog Home](.github/images/blog-home.png) | ![Blog Article](.github/images/blog-article.png) |

| Admin Dashboard | Editor |
| --- | --- |
| ![Admin Dashboard](.github/images/admin-dashboard.png) | ![Admin Editor](.github/images/admin-editor.png) |

## Tech Stack

### Server

- **Language**: [Go 1.25](https://golang.org)
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io)
- **Database**: PostgreSQL
- **Authentication**: JWT, OAuth2, Goth
- **API Docs**: Swagger
- **Scheduled Jobs**: [Cron](https://github.com/robfig/cron)
- **Other Integrations**: User-Agent parsing, Feishu SDK, WeChat Official Account

### Admin

- **Framework**: [Vue 3](https://vuejs.org) + [Vite](https://vitejs.dev)
- **UI Library**: [Element Plus](https://element-plus.org)
- **State Utilities**: [VueUse](https://vueuse.org)
- **Markdown Editor**: CodeMirror 6
- **Charts**: ECharts, echarts-wordcloud
- **Other**: TypeScript, Vue Router, Axios, dayjs, SCSS

### Blog

- **Framework**: [Nuxt 4](https://nuxt.com) - a full-stack framework built on Vue
- **Article Rendering**: markdown-it, Highlight.js, Mermaid
- **Styling**: SCSS
- **SEO**: @nuxtjs/seo, sitemap, Atom feed
- **PWA**: @vite-pwa/nuxt
- **Other**: TypeScript, VueUse, dayjs, Lenis, medium-zoom, APlayer

## Deployment

### Docker Compose Setup (Recommended)

1. Create a `.env` file:

```env
# Database Configuration
DB_PASSWORD=your_database_password

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# Site Configuration
API_URL=https://api.yourdomain.com/api/v1
```

2. Create a `docker-compose.yml` file:

```yaml
services:
  server:
    image: talen8/flec-server:latest
    container_name: flec_server
    restart: unless-stopped
    environment:
      DB_HOST: localhost
      DB_PORT: 5432
      DB_NAME: postgres
      DB_USER: postgres
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - "8080:8080"
    volumes:
      - /srv/flecblog:/app/data
    networks:
      - flec-network

  blog:
    image: talen8/flec-blog:latest
    container_name: flec_blog
    restart: unless-stopped
    environment:
      NUXT_PUBLIC_API_URL: ${API_URL}
    ports:
      - "3000:3000"
    networks:
      - flec-network
    depends_on:
      - server

  admin:
    image: talen8/flec-admin:latest
    container_name: flec_admin
    restart: unless-stopped
    environment:
      API_URL: ${API_URL}
    ports:
      - "4000:4000"
    networks:
      - flec-network
    depends_on:
      - server

networks:
  flec-network:
    driver: bridge
```

<details>
<summary>Full configuration with PostgreSQL included</summary>

```yaml
services:
  postgres:
    image: postgres:15-alpine
    container_name: flec_postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: postgres
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - flec-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  server:
    image: talen8/flec-server:latest
    container_name: flec_server
    restart: unless-stopped
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: postgres
      DB_USER: postgres
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
    ports:
      - "8080:8080"
    volumes:
      - /srv/flecblog:/app/data
    networks:
      - flec-network
    depends_on:
      postgres:
        condition: service_healthy

  blog:
    image: talen8/flec-blog:latest
    container_name: flec_blog
    restart: unless-stopped
    environment:
      NUXT_PUBLIC_API_URL: ${API_URL}
    ports:
      - "3000:3000"
    networks:
      - flec-network
    depends_on:
      - server

  admin:
    image: talen8/flec-admin:latest
    container_name: flec_admin
    restart: unless-stopped
    environment:
      API_URL: ${API_URL}
    ports:
      - "4000:4000"
    networks:
      - flec-network
    depends_on:
      - server

networks:
  flec-network:
    driver: bridge

volumes:
  postgres_data:
```

</details>

3. Start the services:

```bash
docker-compose up -d
```

### Access URLs

| Service | URL |
|------|------|
| Blog | http://localhost:3000 |
| Admin | http://localhost:4000 |
| API Docs | http://localhost:8080/swagger/index.html |

## Run from Source

### Requirements

- Node.js 20+ (admin, blog)
- Go 1.25 (server)
- PostgreSQL 12+ (server)

### Database Setup

Install and configure PostgreSQL first, then create the database and user you want to use.

On first startup, the application will automatically run `pkg/database/sql/init_database.sql` to initialize the database, including table creation and seed data.

> ⚠️ **PostgreSQL 15+ schema permissions**: If you're using PostgreSQL 15 or later and the database user is not `postgres` (the superuser), you'll need to grant schema permissions manually.
>
> ```bash
> sudo -i -u postgres
> psql -U postgres -d <database_name> -c "GRANT CREATE ON SCHEMA public TO <username>;"
> ```
>
> PostgreSQL 15+ restricts `CREATE` privileges on the `public` schema by default, so this step grants the permissions required by the app.

### Server

```bash
cd server
go mod download
cp .env.example .env
# Edit .env and configure your database connection
go run cmd/main.go
```

### Admin

```bash
cd admin
npm install
cp .env.example .env
# Edit .env and set the API base URL
npm run dev
```

### Blog

```bash
cd blog
npm install
cp .env.example .env
# Edit .env and set the API base URL
npm run dev
```

### Configuration

**Server environment variables**

```env
# JWT settings
JWT_SECRET=your_jwt_secret_key

# Server settings
SERVER_PORT=8080
SERVER_ALLOW_ORIGINS=*

# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=your_database_password
```

**Admin environment variables**

```env
VITE_API_URL=https://api.yourdomain.com/api/v1
```

**Blog environment variables**

```env
NUXT_PUBLIC_API_URL=https://api.yourdomain.com/api/v1
```

## Features

### SSR Rendering

The blog frontend runs on Nuxt 4 with SSR enabled, which provides:

- Better SEO, with fully rendered content available to search engines
- Faster first-page load times
- A smoother reading experience

### SEO

FlecBlog includes a complete SEO setup:

- Dynamic `sitemap.xml`
- `robots.txt`
- Atom feed support
- Open Graph metadata
- Structured data

### API Documentation

Once the server is running, you can access the API documentation here:

```
http://localhost:8080/swagger/index.html
```

## Directory Structure

### Server

```text
server/
├── api/              # API definitions
│   ├── middleware/   # Middleware (auth, CORS, logging, rate limiting, RBAC, etc.)
│   ├── router/       # Route definitions
│   └── v1/           # Versioned API v1 handlers
├── cmd/              # Application entry point
│   └── main.go
├── config/           # Configuration management
├── docs/             # Generated Swagger docs
├── internal/         # Internal business logic
│   ├── dto/          # Data transfer objects
│   ├── model/        # Data models
│   ├── repository/   # Data access layer
│   └── service/      # Business services
├── pkg/              # Reusable packages
├── templates/        # Template files
├── Dockerfile
└── go.mod
```

### Admin

```text
admin/
├── src/
│   ├── api/              # API clients
│   ├── assets/           # Static assets
│   ├── components/       # Shared components
│   ├── layouts/          # Page layouts
│   ├── router/           # Route configuration
│   ├── types/            # TypeScript types
│   ├── utils/            # Utility functions
│   ├── views/            # Page views
│   ├── App.vue           # Root component
│   └── main.ts           # Application entry
├── public/               # Public files
├── index.html            # HTML template
├── vite.config.ts        # Vite config
├── nginx.conf            # Nginx config
└── Dockerfile            # Docker setup
```

### Blog

```text
blog/
├── app/                  # Main app directory
│   ├── assets/           # Static assets
│   ├── components/       # Vue components
│   ├── composables/      # Composable utilities
│   ├── layouts/          # Layouts
│   ├── pages/            # File-based routes
│   ├── plugins/          # Nuxt plugins
│   ├── utils/            # Utility helpers
│   └── app.vue           # Root component
├── public/               # Public files
├── server/               # Server-side code
│   ├── plugins/          # Server plugins
│   └── routes/           # API routes
├── types/                # TypeScript types
├── nuxt.config.ts        # Nuxt config
└── Dockerfile            # Docker setup
```

## Philosophy

FlecBlog is not about making things more complicated. It is about making them more complete.

Minimalism does not have to feel empty.
Polish does not have to feel loud.

This project is meant to be a calm, dependable home for content: something you can write in, build on, and slowly shape into a brand of your own.

## Contributing

Issues and pull requests are always welcome.

## License

[MIT License](LICENSE)

## Contact

If you have questions, feel free to reach out:

- Email: [talen2004@163.com](mailto:talen2004@163.com)
- Issues: [GitHub Issues](https://github.com/talen8/FlecBlog/issues)
