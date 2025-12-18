# Zeabur 部署指南

Zeabur 是一个简单易用的云部署平台，支持全自动部署后端、前端及数据库。本指南将指导你如何在 Zeabur 上免费上线本项目。

## 1. 账号准备

1. 访问 [Zeabur 官网](https://zeabur.com/)。
2. 使用 GitHub 账号登录（**无需绑定信用卡**即可使用免费额度）。

## 2. 部署数据库与 Redis

在 Zeabur 控制面板中创建一个新项目，然后添加服务：

### PostgreSQL 部署
1. 点击 **Create Service** -> **Prebuilt Service** -> **PostgreSQL**。
2. 等待服务启动后，点击进入服务详情，记录 **Connection Strings**（连接字符串）或环境变量（Host, Port, User, Password, Database）。

### Redis 部署
1. 点击 **Create Service** -> **Prebuilt Service** -> **Redis**。
2. 等待启动后，记录 Host 和 Port。

## 3. 部署后端 (nutri-baby-server)

1. 点击 **Create Service** -> **Git** -> 选择你的代码仓库。
2. 在弹出窗口中，将 **Root Directory** 设置为 `nutri-baby-server`。
3. 进入 **Environment Variables** 界面，添加项目所需的变量（参考 `config/config.yaml.example`）：
    - `DB_HOST`: PostgreSQL 的 Host。
    - `DB_PORT`: `5432`。
    - `DB_USER`: `postgres`。
    - `DB_PASSWORD`: 你的数据库密码。
    - `DB_NAME`: `nutri_baby`。
    - `REDIS_HOST`: Redis 的 Host。
    - `REDIS_PORT`: `6379`。
    - (其他必要的配置如 `JWT_SECRET` 等)。
4. Zeabur 会根据 `Dockerfile` 自动构建并启动。
5. 在 **Domain** 标签页点击 **Generate Domain**，获得后端的 API URL（例如 `api-v1.zeabur.app`）。

## 4. 部署前端 (nutri-baby-app)

1. 再次点击 **Create Service** -> **Git** -> 选择同一个代码仓库。
2. 将 **Root Directory** 设置为 `nutri-baby-app`。
3. 添加环境变量：
    - `VITE_API_BASE_URL`: 设置为你上一步生成的 **后端 API URL**。
4. 在 **Domain** 标签页点击 **Generate Domain**，获得前端访问地址。

## 5. 验证部署

1. 访问前端地址，检查页面是否正常加载。
2. 尝试登录或提交数据，确认是否能成功连接后端。

## 优势总结
- **不绑卡**：免费计划足够初期开发测试。
- **全自动**：GitHub 推送代码后自动重新构建。
- **自带 SSL**：生成的域名自带 HTTPS。
