#!/bin/bash

echo "🔧 DeepSeek JSON解析问题修复"
echo ""

echo "🐛 问题分析："
echo "❌ 错误: invalid character 'å' looking for beginning of value"
echo "❌ 原因: DeepSeek返回的响应包含非JSON内容或编码问题"
echo "❌ 位置: parseAnalysisResponse JSON解析失败"
echo ""

echo "✅ 解决方案："
echo "1. **强化系统提示** - 明确要求只返回纯JSON"
echo "2. **智能JSON提取** - 从混合内容中提取JSON部分"
echo "3. **详细错误日志** - 记录原始响应和提取的JSON"
echo "4. **容错处理** - 多种JSON提取策略"
echo ""

echo "🔧 技术实现："
echo "**系统提示增强:**"
echo "- 添加 **重要：最终必须只返回纯JSON格式**"
echo "- 强调 **不要包含任何解释文字或其他内容**"
echo ""

echo "**JSON提取算法:**"
echo "1. 括号匹配算法 - 找到完整的 {...} 结构"
echo "2. 正则表达式提取 - 匹配JSON模式"
echo "3. 返回最长匹配 - 选择最完整的JSON"
echo ""

echo "**调试日志:**"
echo "- Debug: 原始AI响应"
echo "- Debug: 提取的JSON"
echo "- Error: JSON解析失败详情"
echo ""

echo "🧪 测试步骤："
echo "1. 启动服务器: ./tmp/nutri-baby-deepseek --config=config/config.yaml"
echo "2. 设置日志级别为DEBUG (在config.yaml中)"
echo "3. 发送分析请求:"
echo ""
echo 'curl -X POST "http://localhost:8080/api/ai-analysis" \'
echo '     -H "Content-Type: application/json" \'
echo '     -H "Authorization: Bearer test-token" \'
echo '     -d "{"'
echo '       "baby_id": 4,'
echo '       "analysis_type": "feeding",'
echo '       "start_date": "2024-11-01",'
echo '       "end_date": "2024-11-12"'
echo '     }"'"'"
echo ""

echo "4. 检查日志输出:"
echo "   - 查看 '原始AI响应' 日志"
echo "   - 查看 '提取的JSON' 日志"
echo "   - 确认JSON解析成功"
echo ""

echo "🔍 预期结果:"
echo "✅ 不再出现 'invalid character' 错误"
echo "✅ 成功提取并解析JSON响应"
echo "✅ 返回有效的分析结果"
echo ""

echo "💡 JSON提取示例:"
echo "**输入:** '根据分析结果，我给出以下建议：{\"score\":85,\"insights\":[...]}'"
echo "**提取:** '{\"score\":85,\"insights\":[...]}'"
echo ""

echo "🚨 如果仍有问题:"
echo "1. 检查DeepSeek API配置"
echo "2. 验证模型是否支持工具调用"
echo "3. 查看完整的错误日志"
echo "4. 考虑调整系统提示"
echo ""

echo "✅ 修复完成！DeepSeek现在应该能正确解析JSON响应了。"
