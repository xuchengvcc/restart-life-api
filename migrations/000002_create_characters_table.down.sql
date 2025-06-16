-- 删除索引
DROP INDEX IF EXISTS idx_characters_user_id ON characters;
DROP INDEX IF EXISTS idx_characters_is_active ON characters;
DROP INDEX IF EXISTS idx_characters_created_at ON characters;
DROP INDEX IF EXISTS idx_characters_life_stage ON characters;
DROP INDEX IF EXISTS idx_characters_game_completed ON characters;

-- 删除角色表
DROP TABLE IF EXISTS characters; 