# 开发文档

## 环境准备

- **Node.js**: >= 16.0.0 (推荐 18.x)
- **npm**: >= 8.0.0
- **微信开发者工具**: 最新版 (如果开发小程序)

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 启动开发服务

#### H5 开发

```bash
npm run dev:h5
```
启动后访问: `http://localhost:5173/`

#### 微信小程序开发

```bash
npm run dev:mp-weixin
```
命令运行后，打开 **微信开发者工具**，导入 `dist/dev/mp-weixin` 目录。

## 项目结构

```
src/
├── api/            # API 接口封装
├── components/     # 公共组件
├── pages/          # 页面文件
├── static/         # 静态资源 (图片等)
├── stores/         # Pinia 状态管理
├── utils/          # 工具函数
├── App.vue         # 根组件
├── main.ts         # 入口文件
├── manifest.json   # uni-app 配置文件
└── pages.json      # 页面及路由配置
```

## 开发规范

### 代码风格
- 使用 ESLint + Prettier 进行代码格式化。
- 提交代码前请确保无 lint 错误。

### 分支管理
- `main`: 主分支，保持稳定。
- `develop`: 开发分支。
- `feature/*`: 功能分支。

## 常见问题

### 1. `sass` 依赖问题
如果遇到 Sass 相关的编译错误，尝试重新安装依赖：
```bash
npm install sass -D
```

### 2. 接口跨域
开发环境 H5 跨域已在 `vite.config.ts` 中配置代理：
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```
确保后端服务已启动在 `8080` 端口。
