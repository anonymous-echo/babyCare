-- 修复换尿布类型不一致问题
-- 将数据库中错误的 'wet' 和 'dirty' 类型修正为 'pee' 和 'poop'

-- 将 wet 类型更新为 pee
UPDATE diaper_records SET type = 'pee' WHERE type = 'wet';

-- 将 dirty 类型更新为 poop
UPDATE diaper_records SET type = 'poop' WHERE type = 'dirty';

-- 验证更新结果
-- SELECT type, COUNT(*) as count FROM diaper_records GROUP BY type;