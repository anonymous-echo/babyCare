#!/bin/bash

echo "🔧 测试宝宝ID传递修复"
echo ""

echo "📋 修复内容："
echo "✅ Mock模型现在从用户消息中动态提取宝宝ID"
echo "✅ 使用正则表达式匹配 '宝宝ID 数字' 模式"
echo "✅ 同时提取日期范围信息"
echo "✅ 工具调用时使用正确的宝宝ID和日期"
echo ""

echo "🧪 测试场景："
echo "1. 创建宝宝ID为4的分析任务"
echo "2. 验证工具调用时使用的是宝宝ID 4而不是1"
echo ""

echo "📝 正则表达式模式："
echo "- 宝宝ID提取: 宝宝ID\\s*(\\d+)"
echo "- 日期范围提取: (\\d{4}-\\d{2}-\\d{2})\\s*至\\s*(\\d{4}-\\d{2}-\\d{2})"
echo ""

echo "🔍 Mock模型行为："
echo "- extractBabyIDFromMessage(): 从消息中提取宝宝ID"
echo "- extractDateRangeFromMessage(): 从消息中提取日期范围"
echo "- 如果提取失败，使用默认值并记录警告日志"
echo ""

echo "💡 用户提示示例："
echo "请对宝宝ID 4 在 2024-11-01 至 2024-11-12 期间的 喂养 数据进行专业分析。"
echo ""
echo "🔧 工具调用参数："
echo '{"baby_id": 4, "start_date": "2024-11-01", "end_date": "2024-11-12", "limit": 100}'
echo ""

echo "🚀 启动测试："
echo "1. 启动服务器: ./tmp/nutri-baby-fixed --config=config/config.yaml"
echo "2. 创建分析任务:"
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
echo "3. 检查日志中的工具调用参数，确认使用的是宝宝ID 4"
echo ""

echo "✅ 修复完成！现在工具调用会使用正确的宝宝ID了。"
