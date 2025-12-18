-- 009_add_relationship_field.sql
-- 为亲友团成员添加与宝宝的关系字段

-- 在 baby_collaborators 表中添加 relationship 字段
ALTER TABLE baby_collaborators ADD COLUMN IF NOT EXISTS relationship VARCHAR(32) DEFAULT '' COMMENT '与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆, 叔叔, 阿姨等';

-- 在 baby_invitations 表中添加 relationship 字段
ALTER TABLE baby_invitations ADD COLUMN IF NOT EXISTS relationship VARCHAR(32) DEFAULT '' COMMENT '与宝宝的关系: 爸爸, 妈妈, 爷爷, 奶奶, 外公, 外婆等';
