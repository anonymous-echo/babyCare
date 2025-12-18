# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

**宝宝喂养日志 (BabyLog+)** - 一个全栈育儿记录系统，帮助新手父母记录和追踪婴幼儿的成长数据。

### 核心特性

- 👶 **去家庭化架构**: 数据以"宝宝"为中心，支持多协作者共同管理单个或多个宝宝
- 📊 **记录管理**: 喂养、睡眠、排泄、成长、疫苗等多维度记录
- 🔔 **智能提醒**: 喂养提醒、疫苗提醒（微信订阅消息）
- 👥 **协作管理**: 通过小程序码邀请协作者，支持角色权限控制

### 技术栈

**前端** (nutri-baby-app):
- uni-app (Vue 3 + TypeScript)
- WotUI  (UI 组件库)
- Vite 5.2.8
- 目标平台：微信小程序（主要）+ 多端支持

**后端** (nutri-baby-server):
- Go 1.25
- Gin Web 框架
- PostgreSQL + GORM
- Redis (缓存)
- DDD 四层架构 + Wire 依赖注入
- 微信 SDK 集成

## 项目结构

```
nutri-baby/
├── nutri-baby-app/          # 前端小程序
│   ├── src/
│   │   ├── pages/          # 18个功能页面
│   │   ├── api/            # 8个 API 模块
│   │   ├── store/          # 5个状态管理模块
│   │   ├── types/          # TypeScript 类型定义
│   │   └── utils/          # 工具库
│   ├── API.md              # API 接口文档 (1241行)
│   └── package.json
├── nutri-baby-server/       # 后端服务
│   ├── cmd/server/         # 应用入口
│   ├── internal/           # DDD 四层架构
│   │   ├── domain/         # 领域层
│   │   │   ├── entity/    # 领域实体
│   │   │   └── repository/ # 仓储接口
│   │   ├── application/    # 应用层 (服务 + DTO)
│   │   ├── infrastructure/ # 基础设施层 (持久化 + 缓存 + 日志)
│   │   └── interface/      # 接口层 (HTTP 处理器 + 路由)
│   ├── pkg/                # 公共库
│   ├── wire/               # Wire 依赖注入
│   ├── config/             # 配置文件
│   ├── migrations/         # 数据库迁移脚本
│   └── Makefile
├── prd.md                  # 产品需求文档
└── CLAUDE.md               # 本文件
```

## 开发命令

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
# 或
cd wire && wire

# 生成 Swagger API 文档 (修改接口注释后执行)
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

# 清理
make clean              # 清理生成文件 (包括 docs/)

# 查看所有命令
make help
```

## 错误处理规范

### 统一错误处理

项目使用 `pkg/errors` 包统一管理所有错误，各层共享相同的错误类型和错误码，避免重复定义和类型转换。

### 错误码规范

错误码定义在 `pkg/errors` 包中，按以下规则分类：

- `0`: 成功
- `1xxx`: 通用错误
- `2xxx`: 服务器错误
- `3xxx`: 业务逻辑错误

常用错误码示例：

```go
const (
    // 成功
    Success ErrorCode = 0

    // 通用错误 1000-1999
    ParamError       ErrorCode = 1001
    Unauthorized     ErrorCode = 1002
    NotFound         ErrorCode = 1003
    Conflict         ErrorCode = 1004
    PermissionDenied ErrorCode = 1005

    // 服务器错误 2000-2999
    InternalError ErrorCode = 2001
    DatabaseError ErrorCode = 2002
    CacheError    ErrorCode = 2003

    // 业务错误 3000-3999
    UserNotFound      ErrorCode = 3001
    InvalidToken      ErrorCode = 3002
    TokenExpired      ErrorCode = 3003
    BabyNotFound      ErrorCode = 3004
    FamilyNotFound    ErrorCode = 3005
    InvalidInvitation ErrorCode = 3006
    RecordNotFound    ErrorCode = 3007
    VaccineNotFound   ErrorCode = 3008
    InvalidVaccineID  ErrorCode = 3009
)
```

### 错误处理流程

1. **存储层 (Repository)**:
   - 捕获底层错误（如 `gorm.ErrRecordNotFound`）
   - 转换为 `pkg/errors` 中定义的错误类型
   - 使用 `errors.Wrap` 添加上下文信息
   ```go
   if errors.Is(err, gorm.ErrRecordNotFound) {
       return nil, errors.New(errors.NotFound, "记录不存在")
   }
   if err != nil {
       return nil, errors.Wrap(errors.DatabaseError, "查询失败", err)
   }
   ```

2. **服务层 (Service)**:
   - 处理业务逻辑错误
   - 使用 `pkg/errors` 中定义的错误码和错误消息
   - 可以包装错误以添加上下文信息
   ```go
   baby, err := s.repo.GetBabyByID(id)
   if err != nil {
       if errors.Is(err, errors.NotFound) {
           return nil, errors.New(errors.BabyNotFound, "未找到宝宝信息")
       }
       return nil, errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
   }
   ```

3. **接口层 (Handler)**:
   - 处理 HTTP 相关的错误
   - 记录错误日志
   - 将错误转换为统一的 API 响应格式
   ```go
   baby, err := service.GetBabyDetail(id, openID)
   if err != nil {
       switch {
       case errors.Is(err, errors.BabyNotFound):
           response.FailWithError(c, errors.ErrBabyNotFound)
       case errors.Is(err, errors.PermissionDenied):
           response.FailWithError(c, errors.ErrPermissionDenied)
       default:
           log.Error("获取宝宝详情失败", 
               zap.String("baby_id", id),
               zap.String("openid", openID),
               zap.Error(err))
           response.FailWithError(c, errors.ErrInternalServer)
       }
       return
   }
   ```

### 最佳实践

1. **错误创建**
   - 使用 `errors.New()` 创建新错误
   - 使用 `errors.Wrap()` 包装底层错误并添加上下文
   - 错误消息应该清晰、具体，便于调试

2. **错误处理**
   - 在可能失败的地方立即处理错误
   - 使用 `errors.Is()` 检查特定错误
   - 在适当的地方添加错误上下文

3. **日志记录**
   - 在接口层记录详细的错误日志
   - 包含请求ID、参数等上下文信息
   - 避免记录敏感信息

4. **API 响应**
   - 返回统一的错误响应格式
   - 对客户端隐藏内部错误细节
   - 提供有意义的错误代码和消息

### 示例：完整的错误处理流程

```go
// 存储层
func (r *babyRepository) FindByID(ctx context.Context, id string) (*entity.Baby, error) {
    var baby entity.Baby
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&baby).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New(errors.NotFound, "记录不存在")
    }
    if err != nil {
        return nil, errors.Wrap(errors.DatabaseError, "查询宝宝信息失败", err)
    }
    return &baby, nil
}

// 服务层
func (s *babyService) GetBabyDetail(ctx context.Context, id, openID string) (*dto.BabyDTO, error) {
    // 检查权限
    hasPermission, err := s.checkPermission(ctx, id, openID)
    if err != nil {
        return nil, err
    }
    if !hasPermission {
        return nil, errors.New(errors.PermissionDenied, "没有权限访问该宝宝信息")
    }

    // 查询宝宝信息
    baby, err := s.babyRepo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, errors.NotFound) {
            return nil, errors.New(errors.BabyNotFound, "未找到宝宝信息")
        }
        return nil, errors.Wrap(err, "获取宝宝详情失败")
    }

    return baby.ToDTO(), nil
}

// 接口层
func (h *BabyHandler) GetBabyDetail(c *gin.Context) {
    babyID := c.Param("id")
    openID, _ := c.Get("openid")

    baby, err := h.babyService.GetBabyDetail(c.Request.Context(), babyID, openID.(string))
    if err != nil {
        switch {
        case errors.Is(err, errors.BabyNotFound):
            response.FailWithError(c, errors.ErrBabyNotFound)
        case errors.Is(err, errors.PermissionDenied):
            response.FailWithError(c, errors.ErrPermissionDenied)
        default:
            h.logger.Error("获取宝宝详情失败", 
                zap.String("baby_id", babyID),
                zap.String("openid", openID),
                zap.Error(err))
            response.FailWithError(c, errors.ErrInternalServer)
        }
        return
    }

    response.SuccessWithData(c, baby)
}
```

### 错误处理流程

1. **存储层 (Repository)**:
   - 捕获底层错误（如 `gorm.ErrRecordNotFound`）
   - 转换为领域错误（如 `errors.NotFound`）
   - 使用 `errors.Wrap` 添加上下文信息

2. **服务层 (Service)**:
   - 处理业务逻辑错误
   - 将底层错误转换为业务相关的错误
   - 添加业务上下文信息

3. **接口层 (Handler)**:
   - 处理 HTTP 相关的错误
   - 记录错误日志
   - 返回统一的错误响应

### 错误码规范

错误码定义在 `pkg/errors` 包中，按以下规则分类：

- `1xxx`: 通用错误
- `2xxx`: 认证授权错误
- `3xxx`: 业务逻辑错误
- `4xxx`: 资源未找到
- `5xxx`: 服务器内部错误

### 最佳实践

1. 使用 `errors.New` 创建新的错误
2. 使用 `errors.Wrap` 包装错误并添加上下文
3. 使用 `errors.Is` 检查错误类型
4. 在服务层处理所有业务相关的错误
5. 在接口层处理所有 HTTP 相关的错误
6. 记录详细的错误日志，包含请求ID和错误堆栈

### 示例代码

```go
// 存储层示例
func (r *babyRepositoryImpl) FindByID(ctx context.Context, babyID string) (*entity.Baby, error) {
    var baby entity.Baby
    err := r.db.WithContext(ctx).
        Where("baby_id = ?", babyID).
        First(&baby).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.New(errors.NotFound, "baby not found")
    }
    if err != nil {
        return nil, errors.Wrap(errors.DatabaseError, "failed to find baby", err)
    }
    return &baby, nil
}

// 服务层示例
func (s *babyService) GetBabyDetail(ctx context.Context, babyID, openID string) (*dto.BabyDTO, error) {
    if err := s.checkPermission(ctx, babyID, openID); err != nil {
        return nil, err
    }

    baby, err := s.babyRepo.FindByID(ctx, babyID)
    if err != nil {
        if errors.Is(err, errors.NotFound) {
            return nil, errors.New(errors.BabyNotFound, "未找到宝宝信息")
        }
        return nil, errors.Wrap(errors.DatabaseError, "获取宝宝信息失败", err)
    }
    
    // 转换为 DTO 并返回
    return baby.ToDTO(), nil
}

// 接口层示例
func (h *BabyHandler) GetBabyDetail(c *gin.Context) {
    babyID := c.Param("id")
    openID, _ := c.Get("openid")

    baby, err := h.babyService.GetBabyDetail(c.Request.Context(), babyID, openID.(string))
    if err != nil {
        switch {
        case errors.Is(err, errors.BabyNotFound):
            response.FailWithError(c, errors.ErrBabyNotFound)
        case errors.Is(err, errors.PermissionDenied):
            response.FailWithError(c, errors.ErrPermissionDenied)
        default:
            h.logger.Error("获取宝宝详情失败", 
                zap.String("baby_id", babyID),
                zap.String("openid", openID.(string)),
                zap.Error(err))
            response.FailWithError(c, errors.ErrInternalServer.WithError(err))
        }
        return
    }

    response.SuccessWithData(c, baby)
}
```

## 核心架构

### 前端架构

#### 去家庭化数据模型

项目已从"家庭中心"架构重构为"宝宝中心"架构：

```typescript
// 核心实体
UserInfo          // 用户信息 (openid, nickName, avatarUrl, defaultBabyId)
BabyProfile       // 宝宝档案 (babyId, name, birthDate, creatorId)
BabyCollaborator  // 宝宝协作者 (openid, role, accessType, expiresAt)
BabyInvitation    // 宝宝邀请码 (inviteCode, babyId, expiresAt)
```

**关键变更**:
- ❌ 已移除: `FamilyInfo`, `FamilyMember`, `Invitation`
- ✅ 新增: `BabyCollaborator`, `BabyInvitation`
- 用户可以同时协作管理多个宝宝
- 每个宝宝独立管理协作者权限

#### 状态管理 (5个模块)

```typescript
// src/store/index.ts
import { useUserStore } from './user'
import { useBabyStore } from './baby'
import { useCollaboratorStore } from './collaborator'
import { useSubscribeStore } from './subscribe'
```

基于 Vue 3 `reactive` 的简化状态管理方案，无 Vuex/Pinia。

#### 页面结构 (18个页面)

参见 [src/pages.json](nutri-baby-app/src/pages.json):

1. **认证与引导** (3个): login, welcome, user
2. **宝宝管理** (5个): list, edit, invite, qrcode, join
3. **记录功能** (4个): feeding, diaper, sleep, growth
4. **数据呈现** (3个): index (首页), timeline, statistics
5. **疫苗管理** (2个): vaccine, vaccine/manage
6. **设置** (1个): settings/subscribe

#### API 调用模块 (8个)

```typescript
// src/api/
auth.ts       // 登录、Token 刷新
baby.ts       // 宝宝 CRUD、协作者管理、邀请码
feeding.ts    // 喂养记录
sleep.ts      // 睡眠记录
diaper.ts     // 换尿布记录
growth.ts     // 成长记录
vaccine.ts    // 疫苗管理
subscribe.ts  // 订阅消息授权
```

#### HTTP 请求封装

```typescript
// src/utils/request.ts
request<T>(config: RequestConfig): Promise<ApiResponse<T>>

// 特性:
// - 自动添加 Bearer Token
// - 401 自动跳转登录
// - 统一错误处理
// - 环境变量配置: VITE_API_BASE_URL
```

#### UI 组件自动导入

```typescript
// vite.config.ts
UniComponents({ resolvers: [NutResolver()] })

// 使用 NutUI 组件无需手动导入
<nut-button type="primary">按钮</nut-button>
```

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

**核心原则**:
- Domain 层定义接口，Infrastructure 层实现
- 依赖倒置，保证领域层独立性
- 使用 Wire 进行编译时依赖注入

#### 领域实体 (去家庭化架构)

```go
// internal/domain/entity/
User              // 用户实体
Baby              // 宝宝实体 (babyId, creatorId, familyGroup)
BabyCollaborator  // 宝宝协作者 (babyId, openid, role, accessType)
BabyInvitation    // 宝宝邀请码 (inviteCode, babyId, expiresAt)

// 记录实体
FeedingRecord     // 喂养记录 (type: breast/bottle/food)
SleepRecord       // 睡眠记录
DiaperRecord      // 换尿布记录
GrowthRecord      // 成长记录

// 疫苗管理
VaccinePlanTemplate   // 疫苗计划模板
BabyVaccinePlan       // 宝宝疫苗计划
VaccineRecord         // 疫苗接种记录
VaccineReminder       // 疫苗提醒

// 订阅消息
SubscribeMessage      // 订阅消息授权
```

#### 仓储模式

```go
// internal/domain/repository/
// 领域层定义接口
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByOpenID(ctx context.Context, openid string) (*entity.User, error)
    // ...
}

// internal/infrastructure/persistence/
// 基础设施层实现接口
type userRepositoryImpl struct {
    db *gorm.DB
}
```

#### Wire 依赖注入

```go
// wire/wire.go
func InitApp(cfg *config.Config) (*App, error) {
    wire.Build(
        // 基础设施层
        logger.NewLogger,
        persistence.NewDatabase,
        persistence.NewRedis,
        wechat.NewClient,

        // 仓储层
        persistence.NewUserRepository,
        persistence.NewBabyRepository,
        // ...

        // 应用服务层
        service.NewAuthService,
        service.NewBabyService,
        // ...

        // HTTP 处理器
        handler.NewAuthHandler,
        // ...

        // 路由和应用
        router.NewRouter,
        NewApp,
    )
    return &App{}, nil
}
```

**重要**: 修改 `wire/wire.go` 后必须运行 `make wire` 重新生成代码。

#### 统一响应格式

```go
// pkg/response/response.go
type Response struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Timestamp int64       `json:"timestamp"`
}

// 使用示例
response.Success(c, data)
response.Error(c, errs.ErrInvalidParam)
```

#### 错误定义

```go
// pkg/errors/errors.go
var (
    ErrInvalidParam   = errors.New(1001, "参数错误")
    ErrUnauthorized   = errors.New(1002, "未授权")
    ErrNotFound       = errors.New(1003, "资源不存在")
    ErrConflict       = errors.New(1004, "数据冲突")
    // ...
)
```

## 关键配置

### 前端配置

**环境变量** (.env):
```bash
VITE_API_BASE_URL=http://localhost:8080
```

**微信小程序配置** ([src/manifest.json](nutri-baby-app/src/manifest.json)):
```json
{
  "mp-weixin": {
    "appid": "wxf47340979046b474"
  }
}
```

### 后端配置

**配置文件** ([config/config.yaml](nutri-baby-server/config/config.yaml)):
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

## 开发注意事项

### 前端开发

1. **新增页面必须在 pages.json 注册**:
```json
{
  "pages": [
    {
      "path": "pages/新页面/index",
      "style": {
        "navigationBarTitleText": "页面标题"
      }
    }
  ]
}
```

2. **使用 TypeScript 类型**:
```typescript
import type { BabyProfile, BabyCollaborator } from '@/types'
```

3. **路径别名**:
```typescript
import { request } from '@/utils/request'
import { useBabyStore } from '@/store'
```

4. **响应式单位**:
使用 `rpx` 作为响应式单位 (750rpx = 屏幕宽度)

### 后端开发

1. **API 开发流程**:
```
创建 DTO → 实现仓储 → 实现服务 → 实现 Handler → 注册路由 → 更新 Wire
```

2. **Wire 使用规范**:
- Provider 函数必须通过参数注入依赖，不要在函数内部创建
- 修改 `wire/wire.go` 后必须运行 `cd wire && wire`
- 不要手动修改 `wire_gen.go`

3. **GORM 最佳实践**:
```go
// ✅ 使用 Preload 避免 N+1 查询
db.Preload("Collaborators").First(&baby, babyId)

// ✅ 使用事务
err := r.db.Transaction(func(tx *gorm.DB) error {
    // ...
    return nil
})

// ✅ 软删除
db.Delete(&baby, babyId) // 自动软删除
```

4. **Context 传递**:
所有 Repository 方法必须接收 `context.Context`

5. **错误处理**:
使用 `pkg/errors` 中定义的业务错误

### 数据库迁移

**位置**: `nutri-baby-server/migrations/`

**现有迁移**:
- 002_vaccine_plan_templates.sql - 疫苗计划模板
- 003_subscribe_message.sql - 订阅消息
- 004_subscribe_message_onetime.sql - 一次性订阅
- 005_feeding_reminder_flag.sql - 喂养提醒标志
- 006_feeding_reminder_interval.sql - 喂养提醒间隔
- 006_feeding_type_field.sql - 喂养类型字段

**执行迁移**:
```bash
make migrate-up
```

## 重要文档

- **API 文档**: Swagger 格式文档，自动生成在 `nutri-baby-server/docs/` 目录
  - `docs/swagger.json` - JSON 格式 API 文档
  - `docs/swagger.yaml` - YAML 格式 API 文档
  - 生成命令: `make swag`
- **产品需求**: [prd.md](prd.md)
- **后端 README**: [nutri-baby-server/README.md](nutri-baby-server/README.md)

## 核心功能状态

### 已完成 ✅

- 用户认证与授权 (微信登录 + JWT)
- 宝宝档案管理 (CRUD)
- 记录功能 (喂养、睡眠、排泄、成长)
- 疫苗管理 (计划、接种记录、提醒)
- 订阅消息 (喂养提醒、疫苗提醒)
- 数据统计和可视化

### 架构特点

- **去家庭化**: 数据以宝宝为中心，支持灵活的协作关系
- **角色权限**: admin (管理员)、editor (编辑者)、viewer (查看者)
- **临时权限**: 支持设置协作者权限过期时间

## 调试技巧

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

## 常见问题

### 前端

**Q: NutUI 组件无法识别**
A: 检查 `vite.config.ts` 中 `UniComponents` 配置是否正确，确保 `NutResolver()` 已配置。

**Q: 页面 404**
A: 检查 `pages.json` 是否已注册该页面路径。

### 后端

**Q: Wire 编译错误**
A: 通常是缺少 Provider 或循环依赖，检查 `wire/wire.go` 中所有依赖是否已声明。

**Q: 数据库连接失败**
A: 检查 `config/config.yaml` 中数据库配置，确保 PostgreSQL 服务已启动。

**Q: Redis 连接失败**
A: 检查 Redis 服务状态，确保端口和密码配置正确。
