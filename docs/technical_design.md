# Nutri-Baby 技术架构与详细设计

## 1. 总体架构
沿用 Clean Architecture (DDD) 分层架构：
- **API Interface**: Gin Handler (HTTP/JSON).
- **Application**: Service (业务逻辑).
- **Domain**: Entity (GORM Model), Repository Interface.
- **Infrastructure**: GORM Implementation, Redis Cache, Cloud Clients.

## 2. 数据库设计变更 (Schema Changes)

### 2.1 喂养记录优化 (Feeding)
- **表名**: `feeding_records`
- **变更**:
    - 新增列 `breast_side` (varchar 10): 枚举 `left`, `right`, `both` (仅母乳有效).
    - 新增列 `food_content` (text): 辅食的具体内容(暂存).

### 2.2 过敏档案 (Allergies)
- **表名**: `baby_allergies` (新增表)
- **字段**:
    - `id` (PK)
    - `baby_id` (FK)
    - `food_name` (varchar 64): 过敏源名称
    - `reaction` (text): 反应描述
    - `severity_level` (int): 1-5 严重程度
    - `created_at`

### 2.3 云存储优化 (COS)
- **方案**: 采用 **前端直传 (Client-side Direct Upload)** 模式。
- **流程**:
    1. 前端请求 `/api/v1/upload/credential` 获取临时密钥 (STS).
    2. 后端调用腾讯云 CAM 接口生成临时密钥.
    3. 前端使用 `cos-wx-sdk-v5` 直接上传文件到 COS Bucket.
    4. 前端将 COS 返回的 Key/URL 传给后端业务接口 (`/api/v1/user/avatar` 等).
- **优势**: 支持分片上传、断点续传（SDK自带）、节省后端带宽。

### 2.4 其他
- **GrowthRecord**: 已包含 `head_circumference`，无需变更。
- **SleepRecord**: 维持现状，分析结果不存入数据库，实时计算或缓存。

## 3. 接口设计 (API)

### 3.1 上传相关
- `GET /v1/upload/credential` -> `{ tmpSecretId, tmpSecretKey, sessionToken, startTime, expiredTime }`

### 3.2 喂养相关
- `POST /v1/feeding-records` (Payload 更新):
  ```json
  {
    "feedingType": "breast",
    "breastSide": "left",
    "duration": 600
  }
  ```

## 4. 目录结构规范
- `docs/`: 项目文档
- `internal/domain/entity/`: 实体定义
- `migrations/`: SQL 迁移脚本 (自动迁移逻辑已集成在 `database.go`)

## 5. 迁移策略
使用 `database.go` 中的 `AutoMigrate` 功能自动变更 Schema。对于新增的必填字段，设置默认值以兼容旧数据。
