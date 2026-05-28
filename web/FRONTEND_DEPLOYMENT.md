# OpenGEO 前端构建部署文档

> 版本：v1.0.0  
> 更新日期：2026-05-26  
> 框架：Next.js 15 + TypeScript + Ant Design

---

## 目录

- [1. 项目概述](#1-项目概述)
- [2. 环境要求](#2-环境要求)
- [3. 开发环境搭建](#3-开发环境搭建)
- [4. 项目配置](#4-项目配置)
- [5. 构建优化](#5-构建优化)
- [6. 部署方案](#6-部署方案)
- [7. CI/CD 配置](#7-cicd-配置)
- [8. 性能优化](#8-性能优化)
- [9. 监控与日志](#9-监控与日志)
- [10. 常见问题](#10-常见问题)

---

## 1. 项目概述

### 1.1 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| React | 19+ | React 框架，支持 RSC/SSR/ISR |
| vite |8+ |构建工具 |
| TypeScript | 5.x | 类型安全 |
| Ant Design | 5.x | UI 组件库 |
| TanStack Query | 5.x | 服务端状态管理 |
| Zustand | 4.x | 客户端状态管理 |
| Tailwind CSS | 4.x | 原子化 CSS |
| ECharts | 5.x | 数据可视化 |

### 1.2 项目结构

```
web/
├── app/                        # Next.js App Router
│   ├── layout.tsx             # 根布局
│   ├── page.tsx               # 首页
│   ├── globals.css            # 全局样式
│   ├── auth/                  # 认证页面
│   ├── dashboard/             # 仪表盘
│   ├── content/               # 内容管理
│   ├── account/               # 账号管理
│   ├── publish/               # 发布管理
│   ├── schedule/              # 调度管理
│   ├── monitor/               # 监测分析
│   └── settings/              # 系统设置
├── components/                # 组件
│   ├── layout/               # 布局组件
│   ├── ui/                   # UI 组件
│   └── providers/            # Provider 组件
├── hooks/                     # React Query Hooks
├── stores/                    # Zustand 状态
├── lib/                       # 工具库
├── types/                     # TypeScript 类型
├── styles/                    # 样式文件
├── public/                    # 静态资源
├── next.config.js             # Next.js 配置
├── tailwind.config.js         # Tailwind 配置
├── tsconfig.json              # TypeScript 配置
├── package.json               # 依赖配置
└── Dockerfile                 # Docker 配置
```

---

## 2. 环境要求

### 2.1 开发环境

- **Node.js**: 18.17.0 或更高版本
- **npm**: 9.0.0 或更高版本（或 yarn/pnpm）
- **Git**: 2.30.0 或更高版本

### 2.2 生产环境

- **Node.js**: 18.x LTS
- **内存**: 至少 512MB
- **CPU**: 至少 1 核

### 2.3 浏览器支持

| 浏览器 | 版本 |
|--------|------|
| Chrome | 最近 2 个版本 |
| Firefox | 最近 2 个版本 |
| Safari | 最近 2 个版本 |
| Edge | 最近 2 个版本 |

---

## 3. 开发环境搭建

### 3.1 安装依赖

```bash
# 进入前端目录
cd web

# 使用 npm 安装依赖
npm install

# 或使用 yarn
yarn install

# 或使用 pnpm
pnpm install
```

### 3.2 环境变量配置

创建 `.env.local` 文件：

```bash
# API 配置
NEXT_PUBLIC_API_URL=http://localhost:8080

# 应用配置
NEXT_PUBLIC_APP_NAME=OpenGEO
NEXT_PUBLIC_APP_VERSION=1.0.0

# 认证配置
NEXT_PUBLIC_JWT_SECRET=your-secret-key

# 其他配置
NEXT_PUBLIC_GA_ID=your-google-analytics-id
```

### 3.3 启动开发服务器

```bash
# 启动开发服务器
npm run dev

# 或使用 yarn
yarn dev

# 或使用 pnpm
pnpm dev
```

访问 http://localhost:3000

### 3.4 开发命令

```bash
# 开发服务器
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint

# 代码格式化
npm run format

# 构建
npm run build

# 启动生产服务器
npm start
```

---

## 4. 项目配置

### 4.1 Next.js 配置

```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  // 启用 React 严格模式
  reactStrictMode: true,
  
  // 启用 SWC 压缩
  swcMinify: true,
  
  // 图片配置
  images: {
    domains: ['localhost', 'your-cdn-domain.com'],
    formats: ['image/avif', 'image/webp'],
  },
  
  // 环境变量
  env: {
    NEXT_PUBLIC_APP_NAME: process.env.NEXT_PUBLIC_APP_NAME,
  },
  
  // Webpack 配置
  webpack: (config, { buildId, dev, isServer, defaultLoaders, webpack }) => {
    // 自定义配置
    return config;
  },
  
  // 重写规则（开发环境代理）
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ];
  },
  
  // 头部配置
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
        ],
      },
    ];
  },
};

module.exports = nextConfig;
```

### 4.2 TypeScript 配置

```json
{
  "compilerOptions": {
    "target": "ES2017",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [{ "name": "next" }],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### 4.3 Tailwind CSS 配置

```javascript
// tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#e6f7ff',
          100: '#bae7ff',
          200: '#91d5ff',
          300: '#69c0ff',
          400: '#40a9ff',
          500: '#1890ff',
          600: '#096dd9',
          700: '#0050b3',
          800: '#003a8c',
          900: '#002766',
        },
      },
    },
  },
  plugins: [],
};
```

### 4.4 PostCSS 配置

```javascript
// postcss.config.js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
```

### 4.5 Biome 配置

```json
{
  "$schema": "https://biomejs.dev/schemas/1.4.0/schema.json",
  "organizeImports": {
    "enabled": true
  },
  "linter": {
    "enabled": true,
    "rules": {
      "recommended": true
    }
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2,
    "lineWidth": 100
  }
}
```

---

## 5. 构建优化

### 5.1 构建命令

```bash
# 生产构建
npm run build

# 分析构建包大小
npm run build:analyze

# 导出静态站点
npm run export
```

### 5.2 构建优化配置

```javascript
// next.config.js
const nextConfig = {
  // 启用压缩
  compress: true,
  
  // 启用静态优化
  output: 'standalone',
  
  // 实验性功能
  experimental: {
    // 启用优化CSS
    optimizeCss: true,
    
    // 启用滚动恢复
    scrollRestoration: true,
  },
  
  // Webpack 优化
  webpack: (config, { dev, isServer }) => {
    // 生产环境优化
    if (!dev && !isServer) {
      // 代码分割
      config.optimization.splitChunks = {
        chunks: 'all',
        cacheGroups: {
          vendor: {
            test: /[\\/]node_modules[\\/]/,
            name: 'vendors',
            chunks: 'all',
          },
        },
      };
    }
    
    return config;
  },
};
```

### 5.3 图片优化

```javascript
// next.config.js
const nextConfig = {
  images: {
    // 图片格式优化
    formats: ['image/avif', 'image/webp'],
    
    // 图片尺寸
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
    
    // 远程图片域名
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.example.com',
      },
    ],
  },
};
```

### 5.4 字体优化

```typescript
// app/layout.tsx
import { Inter } from 'next/font/google';

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  preload: true,
});

export default function RootLayout({ children }) {
  return (
    <html lang="zh-CN" className={inter.className}>
      <body>{children}</body>
    </html>
  );
}
```

---

## 6. 部署方案

### 6.1 Docker 部署

#### Dockerfile

```dockerfile
# Dockerfile

# 阶段1: 依赖安装
FROM node:18-alpine AS deps
WORKDIR /app

# 复制依赖文件
COPY package.json package-lock.json ./

# 安装依赖
RUN npm ci --only=production && \
    npm cache clean --force

# 阶段2: 构建
FROM node:18-alpine AS builder
WORKDIR /app

# 复制依赖
COPY --from=deps /app/node_modules ./node_modules

# 复制源代码
COPY . .

# 设置环境变量
ENV NEXT_TELEMETRY_DISABLED 1
ENV NODE_ENV production

# 构建应用
RUN npm run build

# 阶段3: 生产镜像
FROM node:18-alpine AS runner
WORKDIR /app

# 设置环境变量
ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

# 创建非 root 用户
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# 复制构建产物
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# 设置权限
RUN chown -R nextjs:nodejs /app

# 切换用户
USER nextjs

# 暴露端口
EXPOSE 3000

# 设置环境变量
ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

# 启动命令
CMD ["node", "server.js"]
```

#### docker-compose.yml

```yaml
# docker-compose.yml
version: '3.8'

services:
  frontend:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: opengeo-frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://gateway:8080
      - NODE_ENV=production
    depends_on:
      - gateway
    networks:
      - opengeo-network
    restart: unless-stopped

  gateway:
    build:
      context: ./gateway
      dockerfile: Dockerfile
    container_name: opengeo-gateway
    ports:
      - "8080:8080"
    networks:
      - opengeo-network
    restart: unless-stopped

networks:
  opengeo-network:
    driver: bridge
```

#### 构建和运行

```bash
# 构建镜像
docker build -t opengeo-frontend ./web

# 运行容器
docker run -p 3000:3000 opengeo-frontend

# 使用 docker-compose
docker-compose up -d frontend
```

### 6.2 Nginx 部署

#### Nginx 配置

```nginx
# nginx.conf

upstream frontend {
    server localhost:3000;
}

upstream backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name your-domain.com;
    
    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;
    
    # 静态资源缓存
    location /_next/static/ {
        alias /var/www/opengeo/.next/static/;
        expires 365d;
        access_log off;
        add_header Cache-Control "public, immutable";
    }
    
    # 图片缓存
    location /images/ {
        alias /var/www/opengeo/public/images/;
        expires 30d;
        access_log off;
    }
    
    # API 代理
    location /api/ {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
    
    # 前端应用
    location / {
        proxy_pass http://frontend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: ws: wss: data: blob: 'unsafe-inline'; frame-ancestors 'self';" always;
}
```

#### 部署步骤

```bash
# 1. 构建前端
cd web
npm run build

# 2. 复制文件到服务器
scp -r .next/ user@server:/var/www/opengeo/
scp -r public/ user@server:/var/www/opengeo/
scp package.json user@server:/var/www/opengeo/

# 3. 在服务器上安装依赖
ssh user@server
cd /var/www/opengeo
npm ci --only=production

# 4. 启动应用
npm start

# 5. 配置 Nginx
sudo cp nginx.conf /etc/nginx/sites-available/opengeo
sudo ln -s /etc/nginx/sites-available/opengeo /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 6.3 PM2 部署

#### PM2 配置

```javascript
// ecosystem.config.js
module.exports = {
  apps: [
    {
      name: 'opengeo-frontend',
      script: 'node_modules/.bin/next',
      args: 'start',
      cwd: '/var/www/opengeo',
      instances: 'max',
      exec_mode: 'cluster',
      autorestart: true,
      watch: false,
      max_memory_restart: '1G',
      env: {
        NODE_ENV: 'production',
        PORT: 3000,
      },
    },
  ],
};
```

#### 部署步骤

```bash
# 1. 安装 PM2
npm install -g pm2

# 2. 构建前端
npm run build

# 3. 启动应用
pm2 start ecosystem.config.js

# 4. 保存进程列表
pm2 save

# 5. 设置开机自启
pm2 startup
```

### 6.4 Vercel 部署

#### vercel.json

```json
{
  "version": 2,
  "builds": [
    {
      "src": "package.json",
      "use": "@vercel/next"
    }
  ],
  "routes": [
    {
      "src": "/api/(.*)",
      "dest": "https://your-backend.com/api/$1"
    }
  ],
  "env": {
    "NEXT_PUBLIC_API_URL": "https://your-backend.com"
  }
}
```

#### 部署步骤

```bash
# 1. 安装 Vercel CLI
npm install -g vercel

# 2. 登录
vercel login

# 3. 部署
vercel

# 4. 生产部署
vercel --prod
```

### 6.5 Docker + Kubernetes 部署

#### k8s-deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: opengeo-frontend
  namespace: opengeo
  labels:
    app: opengeo-frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: opengeo-frontend
  template:
    metadata:
      labels:
        app: opengeo-frontend
    spec:
      containers:
      - name: frontend
        image: opengeo-frontend:latest
        ports:
        - containerPort: 3000
        env:
        - name: NODE_ENV
          value: "production"
        - name: NEXT_PUBLIC_API_URL
          value: "http://opengeo-gateway:8080"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /api/health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /api/health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: opengeo-frontend
  namespace: opengeo
spec:
  selector:
    app: opengeo-frontend
  ports:
  - port: 3000
    targetPort: 3000
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: opengeo-frontend
  namespace: opengeo
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: your-domain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: opengeo-frontend
            port:
              number: 3000
```

#### 部署命令

```bash
# 构建镜像
docker build -t opengeo-frontend:latest ./web

# 推送到镜像仓库
docker tag opengeo-frontend:latest your-registry/opengeo-frontend:latest
docker push your-registry/opengeo-frontend:latest

# 部署到 K8s
kubectl apply -f k8s-deployment.yaml

# 查看状态
kubectl get pods -n opengeo
kubectl get services -n opengeo
kubectl get ingress -n opengeo
```

---

## 7. CI/CD 配置

### 7.1 GitHub Actions

```yaml
# .github/workflows/frontend.yml
name: Frontend CI/CD

on:
  push:
    branches: [main, develop]
    paths:
      - 'web/**'
  pull_request:
    branches: [main]
    paths:
      - 'web/**'

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: web/package-lock.json
      
      - name: Install dependencies
        working-directory: ./web
        run: npm ci
      
      - name: Run lint
        working-directory: ./web
        run: npm run lint
      
      - name: Run type check
        working-directory: ./web
        run: npm run type-check

  test:
    name: Test
    runs-on: ubuntu-latest
    needs: lint
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: web/package-lock.json
      
      - name: Install dependencies
        working-directory: ./web
        run: npm ci
      
      - name: Run tests
        working-directory: ./web
        run: npm test

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: web/package-lock.json
      
      - name: Install dependencies
        working-directory: ./web
        run: npm ci
      
      - name: Build
        working-directory: ./web
        run: npm run build
        env:
          NEXT_PUBLIC_API_URL: ${{ secrets.API_URL }}
      
      - name: Upload build artifact
        uses: actions/upload-artifact@v3
        with:
          name: frontend-build
          path: web/.next/

  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    needs: build
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      
      - name: Download build artifact
        uses: actions/download-artifact@v3
        with:
          name: frontend-build
          path: web/.next/
      
      - name: Build Docker image
        working-directory: ./web
        run: |
          docker build -t ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:${{ github.sha }} .
          docker tag ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:${{ github.sha }} ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:latest
      
      - name: Push to Docker Hub
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:${{ github.sha }}
          docker push ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:latest
      
      - name: Deploy to server
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /var/www/opengeo
            docker pull ${{ secrets.DOCKER_REGISTRY }}/opengeo-frontend:latest
            docker-compose up -d frontend
```

### 7.2 GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - lint
  - test
  - build
  - deploy

variables:
  NODE_VERSION: "18"

lint:
  stage: lint
  image: node:${NODE_VERSION}
  script:
    - cd web
    - npm ci
    - npm run lint
    - npm run type-check
  only:
    changes:
      - web/**/*

test:
  stage: test
  image: node:${NODE_VERSION}
  script:
    - cd web
    - npm ci
    - npm test
  only:
    changes:
      - web/**/*

build:
  stage: build
  image: node:${NODE_VERSION}
  script:
    - cd web
    - npm ci
    - npm run build
  artifacts:
    paths:
      - web/.next/
    expire_in: 1 hour
  only:
    changes:
      - web/**/*

deploy:
  stage: deploy
  image: docker:latest
  services:
    - docker:dind
  script:
    - cd web
    - docker build -t $CI_REGISTRY_IMAGE/frontend:$CI_COMMIT_SHA .
    - docker tag $CI_REGISTRY_IMAGE/frontend:$CI_COMMIT_SHA $CI_REGISTRY_IMAGE/frontend:latest
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker push $CI_REGISTRY_IMAGE/frontend:$CI_COMMIT_SHA
    - docker push $CI_REGISTRY_IMAGE/frontend:latest
  only:
    - main
  needs:
    - build
```

---

## 8. 性能优化

### 8.1 代码分割

```typescript
// 动态导入
import dynamic from 'next/dynamic';

const HeavyComponent = dynamic(() => import('@/components/HeavyComponent'), {
  loading: () => <p>Loading...</p>,
  ssr: false, // 禁用 SSR
});
```

### 8.2 图片优化

```tsx
// 使用 Next.js Image 组件
import Image from 'next/image';

export default function MyImage() {
  return (
    <Image
      src="/images/hero.jpg"
      alt="Hero"
      width={800}
      height={600}
      priority // 首屏图片预加载
      placeholder="blur"
      blurDataURL="data:image/jpeg;base64,..."
    />
  );
}
```

### 8.3 字体优化

```tsx
// 使用 next/font
import { Inter, Noto_Sans_SC } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });
const notoSansSC = Noto_Sans_SC({ 
  subsets: ['latin'],
  weight: ['400', '500', '700'],
});
```

### 8.4 缓存策略

```typescript
// API 缓存
import { useQuery } from '@tanstack/react-query';

export function useContents() {
  return useQuery({
    queryKey: ['contents'],
    queryFn: () => api.contents.list(),
    staleTime: 5 * 60 * 1000, // 5 分钟
    cacheTime: 10 * 60 * 1000, // 10 分钟
  });
}
```

### 8.5 预加载

```tsx
// 预加载关键资源
export default function Head() {
  return (
    <head>
      <link rel="preconnect" href="https://api.example.com" />
      <link rel="dns-prefetch" href="https://api.example.com" />
      <link
        rel="preload"
        href="/fonts/inter.woff2"
        as="font"
        type="font/woff2"
        crossOrigin="anonymous"
      />
    </head>
  );
}
```

### 8.6 Lighthouse 优化

```javascript
// next.config.js
const nextConfig = {
  // 启用实验性优化
  experimental: {
    optimizeCss: true,
    optimizePackageImports: ['antd', '@ant-design/icons'],
  },
  
  // 压缩
  compress: true,
  
  // 输出 standalone
  output: 'standalone',
};
```

---

## 9. 监控与日志

### 9.1 错误监控

```typescript
// lib/error-monitoring.ts

// Sentry 配置
import * as Sentry from '@sentry/nextjs';

Sentry.init({
  dsn: process.env.SENTRY_DSN,
  environment: process.env.NODE_ENV,
  tracesSampleRate: 1.0,
});

// 错误边界
export function ErrorBoundary({ children }) {
  return (
    <Sentry.ErrorBoundary fallback={<p>Something went wrong</p>}>
      {children}
    </Sentry.ErrorBoundary>
  );
}
```

### 9.2 性能监控

```typescript
// lib/performance.ts

// Web Vitals 监控
export function reportWebVitals(metric) {
  const { id, name, label, value } = metric;
  
  // 发送到分析服务
  console.log({ id, name, label, value });
  
  // 或发送到自定义端点
  fetch('/api/vitals', {
    method: 'POST',
    body: JSON.stringify({ id, name, label, value }),
  });
}
```

### 9.3 日志记录

```typescript
// lib/logger.ts

const LOG_LEVELS = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3,
};

class Logger {
  level: number;
  
  constructor(level = 'INFO') {
    this.level = LOG_LEVELS[level] || LOG_LEVELS.INFO;
  }
  
  debug(message: string, data?: any) {
    if (this.level <= LOG_LEVELS.DEBUG) {
      console.debug(`[DEBUG] ${message}`, data);
    }
  }
  
  info(message: string, data?: any) {
    if (this.level <= LOG_LEVELS.INFO) {
      console.info(`[INFO] ${message}`, data);
    }
  }
  
  warn(message: string, data?: any) {
    if (this.level <= LOG_LEVELS.WARN) {
      console.warn(`[WARN] ${message}`, data);
    }
  }
  
  error(message: string, error?: Error) {
    if (this.level <= LOG_LEVELS.ERROR) {
      console.error(`[ERROR] ${message}`, error);
    }
  }
}

export const logger = new Logger(process.env.LOG_LEVEL);
```

---

## 10. 常见问题

### 10.1 构建失败

**问题：** `npm run build` 失败

**解决方案：**
```bash
# 清理缓存
rm -rf .next
rm -rf node_modules
npm install

# 检查 TypeScript 错误
npm run type-check

# 检查内存限制
NODE_OPTIONS="--max-old-space-size=4096" npm run build
```

### 10.2 样式问题

**问题：** Tailwind CSS 样式不生效

**解决方案：**
```bash
# 检查 tailwind.config.js 配置
# 确保 content 路径正确
content: [
  './app/**/*.{js,ts,jsx,tsx,mdx}',
  './components/**/*.{js,ts,jsx,tsx,mdx}',
]

# 重新构建 CSS
npm run build
```

### 10.3 API 代理问题

**问题：** 开发环境 API 请求失败

**解决方案：**
```javascript
// next.config.js
async rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*',
    },
  ];
},
```

### 10.4 内存溢出

**问题：** `JavaScript heap out of memory`

**解决方案：**
```bash
# 增加 Node.js 内存限制
export NODE_OPTIONS="--max-old-space-size=4096"

# 或在 package.json 中设置
"scripts": {
  "build": "NODE_OPTIONS='--max-old-space-size=4096' next build"
}
```

### 10.5 部署后白屏

**问题：** 部署后页面白屏

**解决方案：**
```bash
# 检查环境变量
echo $NEXT_PUBLIC_API_URL

# 检查构建输出
ls -la .next/

# 检查服务器日志
pm2 logs opengeo-frontend
```

### 10.6 静态资源 404

**问题：** 静态资源加载失败

**解决方案：**
```javascript
// next.config.js
const nextConfig = {
  assetPrefix: process.env.CDN_URL || '',
  images: {
    domains: ['your-cdn-domain.com'],
  },
};
```

---

## 附录

### A. 环境变量参考

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| NEXT_PUBLIC_API_URL | 是 | http://localhost:8080 | 后端 API 地址 |
| NEXT_PUBLIC_APP_NAME | 否 | OpenGEO | 应用名称 |
| NEXT_PUBLIC_APP_VERSION | 否 | 1.0.0 | 应用版本 |
| NEXT_PUBLIC_GA_ID | 否 | - | Google Analytics ID |
| SENTRY_DSN | 否 | - | Sentry DSN |
| CDN_URL | 否 | - | CDN 域名 |

### B. Makefile 前端命令

```bash
# 前端开发
make frontend-dev          # 启动前端开发服务器
make frontend-build        # 构建前端
make frontend-install      # 安装前端依赖
make frontend-lint         # 前端代码检查
make frontend-test         # 前端测试

# Docker 命令
make frontend-docker       # 构建前端 Docker 镜像
make frontend-docker-run   # 运行前端 Docker 容器
```

### C. 参考链接

- [Next.js 文档](https://nextjs.org/docs)
- [Ant Design 文档](https://ant.design/)
- [Tailwind CSS 文档](https://tailwindcss.com/)
- [TanStack Query 文档](https://tanstack.com/query)
- [Vercel 部署文档](https://vercel.com/docs)

---

> 文档版本：v1.0.0  
> 最后更新：2026-05-26