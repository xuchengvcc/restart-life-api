-- 删除索引
DROP INDEX IF EXISTS idx_character_user_id ON character_tab;
DROP INDEX IF EXISTS idx_character_is_active ON character_tab;
DROP INDEX IF EXISTS idx_character_created_at ON character_tab;
DROP INDEX IF EXISTS idx_character_life_stage ON character_tab;
DROP INDEX IF EXISTS idx_character_game_completed ON character_tab;

-- 删除角色表
DROP TABLE IF EXISTS character_tab; 