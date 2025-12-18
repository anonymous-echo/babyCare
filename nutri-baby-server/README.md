# Nutri Baby Server

宝宝喂养日志服务端 API

## 技术栈

- **框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **缓存**: Redis
- **日志**: Zap
- **依赖注入**: Wire
- **架构**: DDD + Clean Architecture

## 项目结构

```
nutri-baby-server/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go
├── internal/               # 内部应用代码
│   ├── domain/            # 领域层
│   │   ├── entity/        # 实体
│   │   ├── valueobject/   # 值对象
│   │   ├── errors/        # 领域错误定义
│   │   └── repository/    # 仓储接口
│   ├── application/       # 应用层
│   │   ├── dto/           # 数据传输对象
│   │   ├── service/       # 应用服务
│   │   └── assembler/     # 组装器
│   ├── infrastructure/    # 基础设施层
│   │   ├── persistence/   # 持久化
│   │   ├── cache/         # 缓存
│   │   ├── logger/        # 日志
│   │   └── config/        # 配置
│   └── interface/         # 接口层
│       ├── http/          # HTTP处理器
│       ├── middleware/    # 中间件
│       └── router/        # 路由
├── pkg/                   # 公共库
│   ├── errors/            # 错误定义
│   ├── response/          # 响应封装
│   └── utils/             # 工具函数
├── wire/                  # Wire依赖注入
│   └── wire.go
├── config/                # 配置文件
│   └── config.yaml
├── migrations/            # 数据库迁移
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 安装开发工具

```bash
make install-tools
```

安装的工具包括：
- `wire` - 依赖注入代码生成
- `goimports` - Go 代码导入管理
- `golangci-lint` - 代码检查
- `swag` - Swagger API 文档生成

### 3. 生成依赖注入代码

```bash
make wire
```

或者：

```bash
cd wire && wire
```

### 4. 配置数据库

编辑 `config/config.yaml`：

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: your_password
  dbname: nutri_baby
```

### 5. 运行迁移

```bash
make migrate-up
```

### 6. 启动服务

```bash
make run
```

服务启动后，访问 http://localhost:8080/swagger/index.html 查看 API 文档（如果已配置 Swagger UI）

### 7. 生成 Swagger API 文档

修改接口注释后，执行：

```bash
make swag
```

文档会生成到 `docs/` 目录：
- `docs/swagger.json` - JSON 格式
- `docs/swagger.yaml` - YAML 格式

## 开发命令

```bash
# 生成Wire代码
make wire

# 生成Swagger API文档
make swag

# 运行服务
make run

# 运行测试
make test

# 数据库迁移
make migrate-up
make migrate-down

# 代码格式化
make fmt

# 代码检查
make lint

# 清理生成文件（包括 docs/）
make clean

# 查看所有命令
make help
```

## 错误处理最佳实践

### 统一错误处理

项目使用 `pkg/errors` 包统一管理所有错误，各层共享相同的错误类型和错误码，避免重复定义和类型转换。

### 错误处理流程

1. **存储层 (Repository)**:
   - 捕获底层错误（如 `gorm.ErrRecordNotFound`）
   - 转换为 `pkg/errors` 中定义的错误类型
   - 使用 `errors.Wrap` 添加上下文信息
   ```go
   if errors.Is(err, gorm.ErrRecordNotFound) {
       return nil, errors.New(errors.NotFound, "记录不存在")
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
           log.Error("获取宝宝详情失败", zap.Error(err))
           response.FailWithError(c, errors.ErrInternalServer)
       }
       return
   }
   ```

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

### 错误日志

- 在接口层记录详细的错误日志
- 包含请求ID、错误堆栈等调试信息

### 最佳实践

1. 始终使用 `errors.Wrap` 添加上下文信息
2. 在服务层处理所有业务相关的错误
3. 在接口层处理所有 HTTP 相关的错误
4. 使用统一的错误响应格式
5. 记录详细的错误日志，方便问题排查
