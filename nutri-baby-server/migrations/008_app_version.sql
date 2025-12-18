-- 应用版本管理表
-- 创建时间：2025-11-14
-- 功能：管理应用版本信息，支持动态更新版本号、构建信息等

-- 创建版本信息表
CREATE TABLE IF NOT EXISTS app_versions (
    id BIGSERIAL PRIMARY KEY,
    version VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL DEFAULT '宝宝喂养时刻',
    description TEXT,
    min_version VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    force_update BOOLEAN NOT NULL DEFAULT FALSE,
    release_notes TEXT,
    build_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 添加注释
COMMENT ON TABLE app_versions IS '应用版本信息表';
COMMENT ON COLUMN app_versions.id IS '版本ID';
COMMENT ON COLUMN app_versions.version IS '版本号';
COMMENT ON COLUMN app_versions.name IS '应用名称';
COMMENT ON COLUMN app_versions.description IS '版本描述信息';
COMMENT ON COLUMN app_versions.min_version IS '最小支持版本';
COMMENT ON COLUMN app_versions.is_active IS '是否为当前活跃版本';
COMMENT ON COLUMN app_versions.force_update IS '是否强制更新';
COMMENT ON COLUMN app_versions.release_notes IS '发布说明';
COMMENT ON COLUMN app_versions.build_time IS '构建时间';

-- 创建索引
CREATE INDEX idx_app_versions_active ON app_versions(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_app_versions_version ON app_versions(version);

-- 创建自动更新 updated_at 的触发器
CREATE OR REPLACE FUNCTION update_app_version_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_app_version_updated_at ON app_versions;
CREATE TRIGGER update_app_version_updated_at
    BEFORE UPDATE ON app_versions
    FOR EACH ROW
    EXECUTE FUNCTION update_app_version_updated_at();

-- 插入初始版本数据
INSERT INTO app_versions (version, name, description, is_active, build_time)
VALUES (
    '3.2.0',
    '宝宝喂养时刻',
    '修复隐私政策合规问题，优化版本管理系统',
    FALSE,
    CURRENT_TIMESTAMP
) ON CONFLICT (version) DO NOTHING;

-- 插入新版本数据（3.3.0）
INSERT INTO app_versions (version, name, description, is_active, force_update, release_notes, build_time)
VALUES (
    '3.3.0',
    '宝宝喂养时刻',
    '新增应用版本管理功能，支持版本检测和更新提醒',
    TRUE,
    FALSE,
    '✨ 新增功能：\n- 应用版本管理系统\n- 版本检测和更新提醒\n- 登录页面版本检测\n- 统计页面睡眠数据显示优化\n- 时间线和喂养记录页面UI改进',
    CURRENT_TIMESTAMP
) ON CONFLICT (version) DO NOTHING;

-- 创建视图：获取当前活跃版本
CREATE OR REPLACE VIEW current_app_version AS
SELECT
    id,
    version,
    name,
    description,
    min_version,
    force_update,
    release_notes,
    build_time,
    created_at,
    updated_at
FROM app_versions
WHERE is_active = TRUE
ORDER BY created_at DESC
LIMIT 1;

-- 权限设置（根据实际部署需要调整）
-- GRANT SELECT ON app_versions TO app_user;
-- GRANT SELECT ON current_app_version TO app_user;
