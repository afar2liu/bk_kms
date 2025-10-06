# bk_kms 前端项目使用说明

## 项目简介

基于 **Vue 3 + TypeScript + ArcoVue** 的书签知识管理系统前端项目。

## 技术栈

- **框架**: Vue 3.4 + TypeScript
- **UI 组件库**: Arco Design Vue 2.55
- **路由**: Vue Router 4.3
- **状态管理**: Pinia 2.1
- **HTTP 客户端**: Axios 1.6
- **构建工具**: Vite 5.1

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 请求封装
│   │   ├── auth.ts       # 认证相关
│   │   ├── bookmark.ts   # 书签相关
│   │   └── tag.ts        # 标签相关
│   ├── components/       # 公共组件
│   │   ├── ImportModal.vue           # 导入书签对话框（SSE 实时进度）
│   │   ├── BookmarkFormModal.vue     # 书签创建/编辑表单
│   │   └── ContentViewModal.vue      # 书签内容查看
│   ├── layouts/          # 布局组件
│   │   └── MainLayout.vue            # 主布局
│   ├── stores/           # 状态管理
│   │   └── user.ts                   # 用户状态
│   ├── types/            # TypeScript 类型定义
│   │   └── index.ts
│   ├── utils/            # 工具函数
│   │   ├── request.ts                # Axios 封装
│   │   └── storage.ts                # LocalStorage 封装
│   ├── views/            # 页面组件
│   │   ├── Login.vue                 # 登录页
│   │   ├── BookmarkList.vue          # 书签列表页
│   │   └── TagManage.vue             # 标签管理页
│   ├── router/           # 路由配置
│   │   └── index.ts
│   ├── App.vue           # 根组件
│   ├── main.ts           # 入口文件
│   └── env.d.ts          # 类型声明
├── package.json
├── vite.config.ts
├── tsconfig.json
└── index.html
```

## 快速开始

### 1. 安装依赖

```bash
cd frontend
npm install
# 或
pnpm install
# 或
yarn install
```

### 2. 启动开发服务器

```bash
npm run dev
```

访问: http://localhost:3000

### 3. 构建生产版本

```bash
npm run build
```

构建产物在 `dist/` 目录。

### 4. 预览生产构建

```bash
npm run preview
```

## 功能特性

### ✅ 用户认证
- 登录页面（含图形验证码）
- Token 自动管理
- 路由守卫

### ✅ 书签管理
- **列表展示**: 表格形式，支持分页
- **搜索过滤**: 关键词搜索 + 标签筛选
- **创建书签**: 表单创建，支持自动获取网页内容
- **编辑书签**: 修改书签信息
- **删除书签**: 单个删除 + 批量删除
- **查看内容**: 查看归档的网页内容
- **批量导入**: 上传 HTML 文件，SSE 实时进度推送

### ✅ 标签管理
- 标签列表展示（按使用次数排序）
- 标签重命名
- 标签使用统计

### ✅ 书签导入（SSE 实时进度）
- 支持 Netscape Bookmark 格式
- 实时进度条
- 详细日志输出
- 成功/跳过/失败统计
- 可选的自动标签生成

## API 代理配置

开发环境下，Vite 会自动代理 API 请求到后端服务器：

```typescript
// vite.config.ts
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    },
    '/bookmark': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

## 环境变量

创建 `.env.local` 文件（不会被 git 跟踪）：

```bash
# API 基础 URL（可选，默认使用代理）
VITE_API_BASE_URL=http://localhost:8080
```

## 核心实现说明

### 1. SSE 书签导入

```typescript
// src/api/bookmark.ts
export async function importBookmarks(
  file: File,
  generateTag: boolean,
  onProgress: (event: ImportProgressEvent) => void
): Promise<void> {
  const formData = new FormData()
  formData.append('bookmark_file', file)

  const response = await fetch('/api/v1/bookmarks/import?generate_tag=' + generateTag, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` },
    body: formData
  })

  const reader = response.body?.getReader()
  const decoder = new TextDecoder()

  while (true) {
    const { value, done } = await reader.read()
    if (done) break

    const chunk = decoder.decode(value)
    const lines = chunk.split('\n')

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.substring(6))
        onProgress(data)
      }
    }
  }
}
```

### 2. Axios 拦截器

```typescript
// 请求拦截器：自动添加 Token
request.interceptors.request.use((config) => {
  const token = useUserStore().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => {
    if (data.code !== 0) {
      Message.error(data.msg)
      return Promise.reject(new Error(data.msg))
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，跳转登录
      useUserStore().logout()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

### 3. 路由守卫

```typescript
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const requiresAuth = to.matched.some((record) => 
    record.meta.requiresAuth !== false
  )

  if (requiresAuth && !userStore.isLoggedIn()) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})
```

## 常见问题

### 1. 端口冲突

修改 `vite.config.ts` 中的 `server.port`。

### 2. API 请求失败

检查后端服务是否启动（http://localhost:8080）。

### 3. 跨域问题

开发环境已配置代理，生产环境需要后端配置 CORS。

### 4. Token 过期

自动跳转到登录页，重新登录即可。

## 开发建议

1. **代码规范**: 使用 ESLint 检查代码
2. **类型安全**: 充分利用 TypeScript 类型系统
3. **组件复用**: 抽取公共组件，避免重复代码
4. **状态管理**: 合理使用 Pinia，避免 prop drilling
5. **性能优化**: 使用 `v-memo`、`computed` 等优化渲染

## 部署

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/bk_kms/frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /bookmark {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 许可证

MIT
