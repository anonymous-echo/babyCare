# Zeabur 部署指南

Zeabur 是一个简单易用的云部署平台，支持全自动部署后端、前端及数据库。本指南将指导你如何在 Zeabur 上免费上线本项目。

## 1. 账号准备

1. 访问 [Zeabur 官网](https://zeabur.com/)。
2. 使用 GitHub 账号登录（**无需绑定信用卡**即可使用免费额度）。

## 2. 为什么会 404？

在 Zeabur 部署时遇到 404 通常由以下原因引起：
1. **SPA 路由问题**：前端（Vue/Uni-app）使用 `History` 模式时，直接访问非首页路径（如 `/pages/index/index`）并刷新页面，服务器（Nginx/Go）由于找不到对应的物理文件会返回 404。
2. **前后端分离部署复杂性**：分别部署前端和后端时，前端如果无法通过跨域配置或服务发现找到后端，也会导致 API 请求失败。

## 3. 一键部署方案（推荐）

为了简化部署并彻底解决 404 问题，我们提供了一个**全栈统一 Dockerfile**。它将前端静态文件打包进 Go 后程，由 Go Server 直接托管，统一通过一个域名访问。

### 部署步骤

1. 在 Zeabur 中创建一个新项目。
2. 添加 **PostgreSQL** 和 **Redis** 服务并记录连接信息。
3. 点击 **Create Service** -> **Git** -> 选择你的代码仓库。
4. **关键配置**：
    - **Root Directory**: 保持为空（即项目根目录）。
    - Zeabur 会自动识别根目录下的 `Dockerfile`。
5. **添加环境变量**：
    - `DB_HOST`: PostgreSQL 的 Host。
    - `DB_PORT`: `5432`。
    - `DB_USER`: `postgres`。
    - `DB_PASSWORD`: 你的数据库密码。
    - `DB_NAME`: `nutri_baby`。
    - `REDIS_HOST`: Redis 的 Host。
    - `REDIS_PORT`: `6379`。
    - `PORT`: `80`。
    - `JWT_SECRET`: 随机字符串。
    - `GIN_MODE`: `release`。
6. 在 **Domain** 标签页点击 **Generate Domain**。

### 方案优势
- **自动处理 404**：Go 后端已配置 SPA 路由兼容，所有位置路径都会重定向到 `index.html`。
- **无需处理跨域**：前后端同源，API 请求直接发往 `/api/v1`。
- **一键部署**：只需配置一套环境变量。

## 4. 后端配置 (nutri-baby-server)

如果你坚持要分开部署，请确保：
1. 进入 **Environment Variables** 界面，添加项目所需的变量。
2. Zeabur 会根据 `nutri-baby-server/Dockerfile` 自动构建。
3. 使用生成的后端域名更新前端的 `VITE_API_BASE_URL`。

## 5. 验证部署

1. 访问生成的域名，检查页面是否能够正常刷新。
2. 尝试登录和记录数据，确认后端 API 通信正常。
