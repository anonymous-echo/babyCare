# 🍼 BabyLog+ - 宝宝喂养日志系统

一个全栈育儿记录系统，帮助新手父母记录和追踪婴幼儿的成长数据。支持多协作者共同管理，提供喂养提醒、疫苗管理和数据统计等功能。

## 🌟 核心特性

### 👶 去家庭化架构
- 数据以"宝宝"为中心，支持多协作者共同管理
- 灵活的权限控制：管理员、编辑者、查看者
- 支持临时协作者权限设置

### 📊 多维度记录管理
- **喂养记录**：母乳、配方奶、辅食记录
- **睡眠记录**：睡眠时间统计和分析
- **排泄记录**：尿布更换追踪
- **成长记录**：身高体重数据管理
- **疫苗管理**：接种计划和提醒

### 🔔 智能提醒系统
- 喂养提醒（基于订阅消息）
- 疫苗接种提醒
- 支持微信订阅消息推送

### 👥 协作管理
- 通过小程序码邀请协作者
- 角色权限控制
- 协作者通知机制

## 🏗️ 技术架构

### 前端 (nutri-baby-app)
- **框架**: uni-app (Vue 3 + TypeScript)
- **UI库**: WotUI 组件库
- **构建工具**: Vite 5.2.8
- **状态管理**: Vue 3 reactive (无 Vuex/Pinia)
- **目标平台**: 微信小程序（主要）+ 多端支持

### 后端 (nutri-baby-server)
- **语言**: Go 1.25
- **Web框架**: Gin
- **数据库**: PostgreSQL + GORM
- **缓存**: Redis
- **架构**: DDD 四层架构 + Wire 依赖注入
- **文档**: Swagger API 文档

## 📁 项目结构

```
nutri-baby/
├── nutri-baby-app/          # 前端小程序
│   ├── src/
│   │   ├── api/            # 12个 API 模块
│   │   ├── pages/          # 18个功能页面
│   │   ├── store/          # 5个状态管理模块
│   │   ├── types/          # TypeScript 类型定义
│   │   └── utils/          # 工具库
│   ├── API.md              # API 接口文档
│   └── package.json
├── nutri-baby-server/       # 后端服务
│   ├── cmd/server/         # 应用入口
│   ├── internal/           # DDD 四层架构
│   │   ├── domain/         # 领域层
│   │   ├── application/    # 应用层
│   │   ├── infrastructure/ # 基础设施层
│   │   └── interface/      # 接口层
│   ├── pkg/                # 公共库
│   ├── wire/               # Wire 依赖注入
│   ├── config/             # 配置文件
│   ├── migrations/         # 数据库迁移
│   └── Makefile
├── prd.md                   # 产品需求文档
└── CLAUDE.md               # 项目详细说明
```

## 🚀 快速开始

### 环境要求
- Node.js 16+
- Go 1.25+
- PostgreSQL 12+
- Redis 6+
- 微信开发者工具

### 前端开发

```bash
# 进入前端目录
cd nutri-baby-app

# 安装依赖
npm install

# 微信小程序开发 (主要平台)
npm run dev:mp-weixin

# H5 开发
npm run dev:h5

# 类型检查
npm run type-check

# 构建微信小程序
npm run build:mp-weixin
```

### 后端开发

```bash
# 进入后端目录
cd nutri-baby-server

# 安装依赖
go mod download

# 安装开发工具 (首次)
make install-tools

# 生成 Wire 依赖注入代码 (修改 wire.go 后必须执行)
make wire

# 生成 Swagger API 文档
make swag

# 运行服务 (默认端口 8080)
make run

# 构建可执行文件
make build              # 当前操作系统
make build-linux        # Linux amd64
make build-all          # 所有平台

# 测试和代码质量
make test               # 运行测试
make fmt                # 代码格式化
make lint               # 代码检查

# 数据库迁移
make migrate-up         # 执行迁移
make migrate-down       # 回滚迁移
```

## 📋 功能模块

### 页面结构 (18个页面)

1. **认证与引导** (3个)
   - 登录页
   - 欢迎页
   - 用户信息页

2. **宝宝管理** (5个)
   - 宝宝列表
   - 宝宝编辑
   - 邀请协作者
   - 邀请码生成
   - 加入宝宝

3. **记录功能** (4个)
   - 喂养记录
   - 睡眠记录
   - 换尿布记录
   - 成长记录

4. **数据呈现** (3个)
   - 首页统计
   - 时间线
   - 数据统计

5. **疫苗管理** (2个)
   - 疫苗计划
   - 疫苗管理

6. **设置** (1个)
   - 订阅设置

### API 模块 (12个)

- `auth` - 用户认证
- `baby` - 宝宝管理
- `feeding` - 喂养记录
- `sleep` - 睡眠记录
- `diaper` - 换尿布记录
- `growth` - 成长记录
- `vaccine` - 疫苗管理
- `subscribe` - 订阅消息
- `statistics` - 数据统计
- `timeline` - 时间线
- `collaborator` - 协作者管理
- `invitation` - 邀请管理

## 🔧 配置说明

### 前端配置

**环境变量** (.env):
```bash
VITE_API_BASE_URL=http://localhost:8080
```

**微信小程序配置** (src/manifest.json):
```json
{
  "mp-weixin": {
    "appid": "your_wechat_appid"
  }
}
```

### 后端配置

**配置文件** (config/config.yaml):
```yaml
server:
  port: 8080
  mode: debug # debug, release, test

database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: nutri_baby

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key-change-in-production
  expire_hours: 72

wechat:
  app_id: your_wechat_app_id
  app_secret: your_wechat_app_secret
  subscribe_templates:
    breast_feeding_reminder: "TEMPLATE_ID"
    bottle_feeding_reminder: "TEMPLATE_ID"
    vaccine_reminder: "TEMPLATE_ID"
```

## 🏗️ 架构设计

### 前端架构

#### 去家庭化数据模型
项目采用"宝宝中心"架构，用户可以同时协作管理多个宝宝：

```typescript
// 核心实体
UserInfo          // 用户信息
BabyProfile       // 宝宝档案
BabyCollaborator  // 宝宝协作者
BabyInvitation    // 宝宝邀请码
```

#### 状态管理
基于 Vue 3 `reactive` 的简化状态管理方案：
- `useUserStore` - 用户状态
- `useBabyStore` - 宝宝状态
- `useCollaboratorStore` - 协作者状态
- `useSubscribeStore` - 订阅状态
- `useRecordStore` - 记录状态

### 后端架构

#### DDD 四层架构
```
Interface Layer (接口层)
  ↓ 依赖
Application Layer (应用层)
  ↓ 依赖
Domain Layer (领域层) ← Infrastructure Layer (基础设施层)
                        ↑ 实现仓储接口
```

#### 领域实体
```go
// 用户相关
User              // 用户实体
Baby              // 宝宝实体
BabyCollaborator  // 协作者
BabyInvitation    // 邀请码

// 记录相关
FeedingRecord     // 喂养记录
SleepRecord       // 睡眠记录
DiaperRecord      // 换尿布记录
GrowthRecord      // 成长记录

// 疫苗管理
VaccinePlanTemplate   // 疫苗计划模板
BabyVaccinePlan       // 宝宝疫苗计划
VaccineRecord         // 疫苗接种记录
VaccineReminder       // 疫苗提醒
```

## 🛠️ 开发规范

### 错误处理

项目使用统一的错误处理机制，错误码分类如下：

- `0`: 成功
- `1xxx`: 通用错误
- `2xxx`: 服务器错误
- `3xxx`: 业务逻辑错误

### 代码规范

1. **前端开发**
   - 使用 TypeScript 类型
   - 路径别名：`@/` 指向 src 目录
   - 响应式单位：使用 `rpx`
   - 新增页面必须在 `pages.json` 注册

2. **后端开发**
   - 所有 Repository 方法必须接收 `context.Context`
   - 使用 Wire 进行依赖注入
   - 统一使用 `pkg/errors` 处理错误
   - API 响应使用统一格式

### 数据库迁移

迁移文件位于 `migrations/` 目录，使用以下命令：
```bash
make migrate-up    # 执行迁移
make migrate-down  # 回滚迁移
```

## 📚 文档资源

- **[产品需求文档](prd.md)** - 详细的产品功能需求
- **[API 文档](nutri-baby-app/API.md)** - 前端 API 接口文档 (1241行)
- **[后端 README](nutri-baby-server/README.md)** - 后端详细技术文档
- **[Swagger API](nutri-baby-server/docs/swagger.yaml)** - RESTful API 文档
- **[项目说明](CLAUDE.md)** - 完整的项目开发指南

## 🔍 调试技巧

### 前端调试
1. **微信开发者工具**: 查看 Console、Network、Storage
2. **类型检查**: `npm run type-check`
3. **查看编译输出**: `nutri-baby-app/dist/dev/mp-weixin/`

### 后端调试
1. **查看日志**: `nutri-baby-server/logs/app.log`
2. **数据库查询**:
   ```bash
   psql -h localhost -U postgres -d nutri_baby
   ```
3. **Redis 调试**:
   ```bash
   redis-cli -h localhost -p 6379
   ```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [uni-app](https://uniapp.dcloud.io/) - 跨平台开发框架
- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [GORM](https://gorm.io/) - Go ORM 库
- [Wire](https://github.com/google/wire) - Go 依赖注入工具

---

**BabyLog+** - 让每个成长瞬间都值得记录 📸✨