# 部署文档

## 概述

本项目前端采用 uni-app 开发，支持多端发布。本文档主要介绍 H5 和微信小程序的部署流程。

## H5 部署

### 1. 编译构建

在项目根目录下运行以下命令进行构建：

```bash
# 构建 H5 生产环境版本
npm run build:h5
```

构建完成后，生成的静态资源位于 `dist/build/h5` 目录下。

### 2. Nginx 配置

推荐使用 Nginx 部署 H5 应用。

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态资源
    location / {
        root /path/to/your/project/dist/build/h5;
        try_files $uri $uri/ /index.html;
        index index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### 3. Docker 部署 (可选)

如果希望将前端也容器化，可以使用以下 `Dockerfile`:

```dockerfile
FROM nginx:alpine
COPY dist/build/h5 /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 微信小程序部署

### 1. 编译构建

```bash
# 构建微信小程序
npm run build:mp-weixin
```

构建产物位于 `dist/build/mp-weixin`。

### 2. 上传发布

1. 打开 **微信开发者工具**。
2. 导入项目，目录选择 `dist/build/mp-weixin`。
3. 填写正确的 AppID。
4. 测试无误后，点击 **上传** 按钮上传代码至微信公众平台。
5. 在微信公众平台后台提交审核并发布。

## Docker 全栈部署 (推荐)

本项目支持使用 Docker Compose 一键启动前后端全栈服务。

### 1. 准备工作

在项目根目录下创建一个 `.env` 文件，参考以下配置：

```env
# 数据库配置
POSTGRES_PASSWORD=your_secure_password

# 可以添加其他环境变量
```

### 2. 启动服务

在项目根目录下运行：

```bash
docker compose up -d
```

该命令将：
1. 构建前端镜像（基于 Nginx）。
2. 构建后端镜像（基于 Go）。
3. 启动 Postgres 数据库。
4. 启动 Redis 缓存。

### 3. 访问

- **前端**: `http://localhost`
- **后端 API**: `http://localhost/api/` (由前端 Nginx 转发)
- **直接访问后端**: `http://localhost:8080` (取决于 `docker-compose.yml` 是否暴露)

## 注意事项

- **跨域问题**: H5 部署时需注意 API 跨域问题，通过 Nginx 代理或后端 CORS 配置解决。
- **HTTPS**: 生产环境建议配置 SSL 证书，强制使用 HTTPS。
- **环境配置**: 确保编译时 `.env` 文件中的 `VITE_API_BASE_URL` 配置正确。
- **Docker 秘钥**: 请务必修改 `docker-compose.yml` 或 `.env` 中的默认密码。
