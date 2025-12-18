-- AI分析相关表结构
-- 创建时间：2024年
-- 功能：支持大模型AI分析功能

-- AI分析结果表
CREATE TABLE IF NOT EXISTS ai_analyses (
    id BIGSERIAL PRIMARY KEY,
    baby_id BIGINT NOT NULL REFERENCES babies(id) ON DELETE CASCADE,
    analysis_type VARCHAR(20) NOT NULL CHECK (analysis_type IN ('feeding', 'sleep', 'growth', 'health', 'behavior')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'analyzing', 'completed', 'failed')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    input_data TEXT, -- JSON格式存储输入数据
    result TEXT,     -- JSON格式存储分析结果
    score NUMERIC(3,2) CHECK (score >= 0 AND score <= 100), -- 评分0-100
    insights TEXT[], -- 洞察建议数组
    alerts TEXT[],   -- 异常警告数组
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT ai_analyses_date_range CHECK (end_date >= start_date)
);

-- 创建索引
CREATE INDEX idx_ai_analyses_baby_id ON ai_analyses(baby_id);
CREATE INDEX idx_ai_analyses_analysis_type ON ai_analyses(analysis_type);
CREATE INDEX idx_ai_analyses_status ON ai_analyses(status);
CREATE INDEX idx_ai_analyses_created_at ON ai_analyses(created_at);
CREATE INDEX idx_ai_analyses_baby_type_status ON ai_analyses(baby_id, analysis_type, status);

-- 每日建议表
CREATE TABLE IF NOT EXISTS daily_tips (
    id BIGSERIAL PRIMARY KEY,
    baby_id BIGINT NOT NULL REFERENCES babies(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    tips TEXT NOT NULL, -- JSON格式存储建议数组
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT daily_tips_unique_baby_date UNIQUE (baby_id, date)
);

-- 创建索引
CREATE INDEX idx_daily_tips_baby_id ON daily_tips(baby_id);
CREATE INDEX idx_daily_tips_date ON daily_tips(date);
CREATE INDEX idx_daily_tips_expired_at ON daily_tips(expired_at);

-- 创建更新时间的触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为ai_analyses表创建更新时间触发器
DROP TRIGGER IF EXISTS update_ai_analyses_updated_at ON ai_analyses;
CREATE TRIGGER update_ai_analyses_updated_at
    BEFORE UPDATE ON ai_analyses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加注释
COMMENT ON TABLE ai_analyses IS 'AI分析结果表';
COMMENT ON COLUMN ai_analyses.id IS '分析ID';
COMMENT ON COLUMN ai_analyses.baby_id IS '宝宝ID';
COMMENT ON COLUMN ai_analyses.analysis_type IS '分析类型：feeding(喂养), sleep(睡眠), growth(成长), health(健康), behavior(行为)';
COMMENT ON COLUMN ai_analyses.status IS '分析状态：pending(待分析), analyzing(分析中), completed(已完成), failed(失败)';
COMMENT ON COLUMN ai_analyses.start_date IS '分析开始日期';
COMMENT ON COLUMN ai_analyses.end_date IS '分析结束日期';
COMMENT ON COLUMN ai_analyses.input_data IS '输入数据(JSON格式)';
COMMENT ON COLUMN ai_analyses.result IS '分析结果(JSON格式)';
COMMENT ON COLUMN ai_analyses.score IS '综合评分(0-100)';
COMMENT ON COLUMN ai_analyses.insights IS '洞察建议数组';
COMMENT ON COLUMN ai_analyses.alerts IS '异常警告数组';

COMMENT ON TABLE daily_tips IS '每日建议表';
COMMENT ON COLUMN daily_tips.id IS '建议ID';
COMMENT ON COLUMN daily_tips.baby_id IS '宝宝ID';
COMMENT ON COLUMN daily_tips.date IS '建议日期';
COMMENT ON COLUMN daily_tips.tips IS '建议内容(JSON格式)';
COMMENT ON COLUMN daily_tips.expired_at IS '过期时间';

-- 插入一些示例数据（可选）
-- 注意：以下数据仅为示例，实际使用时应该通过API接口生成

-- 创建分析状态统计视图
CREATE OR REPLACE VIEW ai_analysis_stats AS
SELECT
    baby_id,
    analysis_type,
    COUNT(*) as total_count,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_count,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count,
    COUNT(CASE WHEN status IN ('pending', 'analyzing') THEN 1 END) as pending_count,
    AVG(score) as average_score,
    MAX(created_at) as last_analysis_at
FROM ai_analyses
GROUP BY baby_id, analysis_type;

-- 创建每日建议清理函数
CREATE OR REPLACE FUNCTION cleanup_expired_daily_tips()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM daily_tips
    WHERE expired_at < CURRENT_TIMESTAMP;

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- 示例：创建定时任务清理过期数据（需要pg_cron扩展）
-- SELECT cron.schedule('cleanup-expired-tips', '0 2 * * *', 'SELECT cleanup_expired_daily_tips();');

-- 权限设置（根据实际部署需要调整）
-- GRANT SELECT, INSERT, UPDATE ON ai_analyses TO app_user;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON daily_tips TO app_user;
-- GRANT SELECT ON ai_analysis_stats TO app_user;