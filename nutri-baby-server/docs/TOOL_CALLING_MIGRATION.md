# Eino 工具调用架构迁移指南

## 概述

本文档描述了从传统的数据推送模式迁移到 Eino 工具调用架构的完整过程。

## 架构对比

### 旧架构（数据推送）
```
数据收集 → 拼接Prompt → 调用大模型 → 解析结果
```

### 新架构（工具调用）
```
发送任务 → 大模型调用工具 → 获取数据 → 智能分析 → 返回结果
```

## 已实现的组件

### 1. 工具集 (`DataQueryTools`)
- **位置**: `internal/infrastructure/eino/tools/data_query_tools.go`
- **功能**: 提供结构化的数据查询工具
- **工具列表**:
  - `get_baby_info`: 获取宝宝基本信息
  - `get_feeding_data`: 获取喂养记录
  - `get_sleep_data`: 获取睡眠记录
  - `get_growth_data`: 获取成长记录
  - `get_diaper_data`: 获取尿布记录
  - `get_vaccine_data`: 获取疫苗记录

### 2. 增强分析链 (`EnhancedAnalysisChainBuilder`)
- **位置**: `internal/infrastructure/eino/chain/enhanced_analysis_chain.go`
- **功能**: 支持工具调用的智能分析链
- **特性**:
  - 对话循环处理工具调用
  - 动态数据获取
  - 智能分析流程

### 3. 工具调用模型 (`ToolCallingMockChatModel`)
- **位置**: `internal/infrastructure/eino/chain/tool_calling_mock.go`
- **功能**: 模拟工具调用行为
- **用途**: 开发测试和演示

### 4. 增强服务层
- **服务**: `EnhancedAIAnalysisService`
- **处理器**: `EnhancedAIAnalysisHandler`
- **端点**: `/api/ai/enhanced/*`

## 迁移步骤

### 阶段 1: 并行运行（当前状态）
- ✅ 新组件已部署
- ✅ 旧组件保持运行
- ✅ 新端点可用于测试

### 阶段 2: 逐步切换
1. **测试验证**
   ```bash
   # 运行测试脚本
   ./scripts/test_tool_calling.sh
   ```

2. **配置切换**
   ```yaml
   # config.yaml
   ai:
     provider: "mock" # 使用工具调用 Mock
   ```

3. **端点切换**
   - 旧端点: `/api/ai/analysis`
   - 新端点: `/api/ai/enhanced/analysis`

### 阶段 3: 完全迁移
1. 更新前端调用新端点
2. 移除旧组件依赖
3. 清理旧代码

## 使用示例

### 1. 工具调用分析
```bash
curl -X POST "http://localhost:8080/api/ai/enhanced/analysis" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "baby_id": 1,
    "analysis_type": "feeding",
    "start_date": "2024-11-01",
    "end_date": "2024-11-12"
  }'
```

### 2. 工具调用建议生成
```bash
curl -X POST "http://localhost:8080/api/ai/enhanced/daily-tips?baby_id=1&date=2024-11-12" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. 测试工具调用
```bash
curl -X GET "http://localhost:8080/api/ai/enhanced/test-tools?baby_id=1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 配置说明

### AI 提供商配置
```yaml
ai:
  provider: "mock"  # 推荐用于测试工具调用
  # provider: "openai"  # 需要等待 eino-ext 更新
  # provider: "deepseek" # 需要等待 eino-ext 更新
```

### 工具调用流程
1. 用户发送分析请求
2. 大模型接收任务描述
3. 大模型主动调用数据工具
4. 工具返回结构化数据
5. 大模型基于数据生成分析
6. 返回最终结果

## 优势对比

### 旧架构问题
- ❌ 数据冗余加载
- ❌ 硬编码分析逻辑
- ❌ 缺乏灵活性
- ❌ 性能问题

### 新架构优势
- ✅ 按需数据获取
- ✅ 智能分析流程
- ✅ 高度可扩展
- ✅ 符合 Eino 设计理念

## 监控和调试

### 日志关键字
- `ToolCallingMockChatModel`: 工具调用模型日志
- `EnhancedAnalysisChainBuilder`: 增强分析链日志
- `DataQueryTools`: 工具执行日志

### 性能指标
- 工具调用次数
- 数据获取时间
- 分析完成时间
- 错误率

## 故障排除

### 常见问题
1. **工具调用失败**
   - 检查数据库连接
   - 验证宝宝ID有效性
   - 查看工具执行日志

2. **分析超时**
   - 检查 `maxIterations` 设置
   - 优化工具响应时间
   - 调整超时配置

3. **结果解析错误**
   - 验证 JSON 格式
   - 检查字段映射
   - 查看模型输出日志

## 下一步计划

1. **真实模型集成**
   - 等待 eino-ext 更新支持 ToolCallingChatModel
   - 集成 OpenAI/DeepSeek 工具调用

2. **性能优化**
   - 工具调用缓存
   - 并发工具执行
   - 结果缓存策略

3. **功能扩展**
   - 更多数据工具
   - 复杂分析场景
   - 多轮对话支持

## 联系信息

如有问题，请查看：
- 代码注释和文档
- 测试脚本输出
- 日志文件内容
