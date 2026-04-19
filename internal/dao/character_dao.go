package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xuchengvcc/restart-life-api/internal/models"
)

// CharacterDAO 角色数据访问对象接口
type CharacterDAO interface {
	Insert(ctx context.Context, character *models.Character) error
	SelectByID(ctx context.Context, characterID string) (*models.Character, error)
	SelectByUserID(ctx context.Context, userID uint) ([]*models.Character, error)
	SelectActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error
	Delete(ctx context.Context, characterID string) error
	CountByUserIDAndCharacterID(ctx context.Context, userID uint, characterID string) (int, error)
	InsertGameEvent(ctx context.Context, event *models.Event) error
	SelectGameEventsByCharacterID(ctx context.Context, characterID string) ([]models.Event, error)
	UpsertGameDecision(ctx context.Context, decision *models.Decision) error
	SelectGameDecisionByCharacterID(ctx context.Context, characterID string) (*models.Decision, error)
	DeleteGameDecisionByCharacterID(ctx context.Context, characterID string) error
}

// characterDAO MySQL角色数据访问对象实现
type characterDAO struct {
	db *sql.DB
}

// NewCharacterDAO 创建角色数据访问对象
func NewCharacterDAO(db *sql.DB) CharacterDAO {
	return &characterDAO{db: db}
}

// Insert 插入角色
func (d *characterDAO) Insert(ctx context.Context, character *models.Character) error {
	query := `
		INSERT INTO character_tab (
			character_id, user_id, character_name, birth_country, birth_year, current_age,
			gender, race, is_active, created_at, updated_at,
			intelligence, emotional_intelligence, memory, imagination, physical_fitness, appearance,
			education_level, marital_status, current_country, current_location, current_activity,
			personality, career, skill_tendency, family_background, social_relationships,
			career_desc, education_desc,
			life_stage, current_status, happiness_level, health_level, money,
			total_playtime, game_completed, final_age, death_cause
		) VALUES (
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?
		)
	`

	_, err := d.db.ExecContext(ctx, query,
		character.CharacterID, character.UserID, character.CharacterName,
		character.BirthCountry, character.BirthYear, character.CurrentAge,
		character.Gender, character.Race, character.IsActive,
		character.CreatedAt, character.UpdatedAt,
		character.Attributes.Intelligence, character.Attributes.EmotionalIntelligence,
		character.Attributes.Memory, character.Attributes.Imagination,
		character.Attributes.PhysicalFitness, character.Attributes.Appearance,
		character.EducationLevel, character.MaritalStatus, character.CurrentCountry,
		character.CurrentLocation, character.CurrentActivity, character.Personality,
		character.Career, character.SkillTendency, character.FamilyBackground,
		character.SocialRelationships, character.CareerDesc, character.EducationDesc,
		character.LifeStage, character.CurrentStatus, character.HappinessLevel,
		character.HealthLevel, character.Money, character.TotalPlaytime,
		character.GameCompleted, character.FinalAge, character.DeathCause,
	)

	return err
}

// SelectByID 根据ID查询角色
func (d *characterDAO) SelectByID(ctx context.Context, characterID string) (*models.Character, error) {
	query := `
		SELECT
			character_id, user_id, character_name, birth_country, birth_year, current_age,
			gender, race, is_active, created_at, updated_at,
			intelligence, emotional_intelligence, memory, imagination, physical_fitness, appearance,
			education_level, marital_status, current_country, current_location, current_activity,
			personality, career, skill_tendency, family_background, social_relationships,
			career_desc, education_desc,
			life_stage, current_status, happiness_level, health_level, money,
			total_playtime, game_completed, final_age, death_cause
		FROM character_tab
		WHERE character_id = ? AND is_active = 1
	`

	row := d.db.QueryRowContext(ctx, query, characterID)
	character := &models.Character{}

	err := row.Scan(
		&character.CharacterID, &character.UserID, &character.CharacterName,
		&character.BirthCountry, &character.BirthYear, &character.CurrentAge,
		&character.Gender, &character.Race, &character.IsActive,
		&character.CreatedAt, &character.UpdatedAt,
		&character.Attributes.Intelligence, &character.Attributes.EmotionalIntelligence,
		&character.Attributes.Memory, &character.Attributes.Imagination,
		&character.Attributes.PhysicalFitness, &character.Attributes.Appearance,
		&character.EducationLevel, &character.MaritalStatus, &character.CurrentCountry,
		&character.CurrentLocation, &character.CurrentActivity, &character.Personality,
		&character.Career, &character.SkillTendency, &character.FamilyBackground,
		&character.SocialRelationships, &character.CareerDesc, &character.EducationDesc,
		&character.LifeStage, &character.CurrentStatus, &character.HappinessLevel,
		&character.HealthLevel, &character.Money, &character.TotalPlaytime,
		&character.GameCompleted, &character.FinalAge, &character.DeathCause,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, err
	}

	return character, nil
}

// SelectByUserID 根据用户ID查询角色列表
func (d *characterDAO) SelectByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	query := `
		SELECT
			character_id, user_id, character_name, birth_country, birth_year, current_age,
			gender, race, is_active, created_at, updated_at,
			intelligence, emotional_intelligence, memory, imagination, physical_fitness, appearance,
			education_level, marital_status, current_country, current_location, current_activity,
			personality, career, skill_tendency, family_background, social_relationships,
			career_desc, education_desc,
			life_stage, current_status, happiness_level, health_level, money,
			total_playtime, game_completed, final_age, death_cause
		FROM character_tab
		WHERE user_id = ? AND is_active = 1
		ORDER BY created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character

	for rows.Next() {
		character := &models.Character{}
		err := rows.Scan(
			&character.CharacterID, &character.UserID, &character.CharacterName,
			&character.BirthCountry, &character.BirthYear, &character.CurrentAge,
			&character.Gender, &character.Race, &character.IsActive,
			&character.CreatedAt, &character.UpdatedAt,
			&character.Attributes.Intelligence, &character.Attributes.EmotionalIntelligence,
			&character.Attributes.Memory, &character.Attributes.Imagination,
			&character.Attributes.PhysicalFitness, &character.Attributes.Appearance,
			&character.EducationLevel, &character.MaritalStatus, &character.CurrentCountry,
			&character.CurrentLocation, &character.CurrentActivity, &character.Personality,
			&character.Career, &character.SkillTendency, &character.FamilyBackground,
			&character.SocialRelationships, &character.CareerDesc, &character.EducationDesc,
			&character.LifeStage, &character.CurrentStatus, &character.HappinessLevel,
			&character.HealthLevel, &character.Money, &character.TotalPlaytime,
			&character.GameCompleted, &character.FinalAge, &character.DeathCause,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, rows.Err()
}

// SelectActiveByUserID 根据用户ID查询活跃角色列表
func (d *characterDAO) SelectActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	query := `
		SELECT
			character_id, user_id, character_name, birth_country, birth_year, current_age,
			gender, race, is_active, created_at, updated_at,
			intelligence, emotional_intelligence, memory, imagination, physical_fitness, appearance,
			education_level, marital_status, current_country, current_location, current_activity,
			personality, career, skill_tendency, family_background, social_relationships,
			career_desc, education_desc,
			life_stage, current_status, happiness_level, health_level, money,
			total_playtime, game_completed, final_age, death_cause
		FROM character_tab
		WHERE user_id = ? AND is_active = 1 AND game_completed = 0
		ORDER BY created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*models.Character

	for rows.Next() {
		character := &models.Character{}
		err := rows.Scan(
			&character.CharacterID, &character.UserID, &character.CharacterName,
			&character.BirthCountry, &character.BirthYear, &character.CurrentAge,
			&character.Gender, &character.Race, &character.IsActive,
			&character.CreatedAt, &character.UpdatedAt,
			&character.Attributes.Intelligence, &character.Attributes.EmotionalIntelligence,
			&character.Attributes.Memory, &character.Attributes.Imagination,
			&character.Attributes.PhysicalFitness, &character.Attributes.Appearance,
			&character.EducationLevel, &character.MaritalStatus, &character.CurrentCountry,
			&character.CurrentLocation, &character.CurrentActivity, &character.Personality,
			&character.Career, &character.SkillTendency, &character.FamilyBackground,
			&character.SocialRelationships, &character.CareerDesc, &character.EducationDesc,
			&character.LifeStage, &character.CurrentStatus, &character.HappinessLevel,
			&character.HealthLevel, &character.Money, &character.TotalPlaytime,
			&character.GameCompleted, &character.FinalAge, &character.DeathCause,
		)
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, rows.Err()
}

// Update 更新角色信息
func (d *characterDAO) Update(ctx context.Context, character *models.Character) error {
	query := `
		UPDATE character_tab SET
			character_name = ?, education_level = ?, marital_status = ?, current_country = ?,
			current_location = ?, current_activity = ?, personality = ?, career = ?,
			skill_tendency = ?, family_background = ?, social_relationships = ?,
			career_desc = ?, education_desc = ?, life_stage = ?, current_status = ?,
			happiness_level = ?, health_level = ?, money = ?, total_playtime = ?,
			game_completed = ?, final_age = ?, death_cause = ?, updated_at = ?
		WHERE character_id = ? AND is_active = 1
	`

	_, err := d.db.ExecContext(ctx, query,
		character.CharacterName, character.EducationLevel, character.MaritalStatus,
		character.CurrentCountry, character.CurrentLocation, character.CurrentActivity,
		character.Personality, character.Career, character.SkillTendency,
		character.FamilyBackground, character.SocialRelationships, character.CareerDesc,
		character.EducationDesc, character.LifeStage, character.CurrentStatus,
		character.HappinessLevel, character.HealthLevel, character.Money,
		character.TotalPlaytime, character.GameCompleted, character.FinalAge,
		character.DeathCause, character.UpdatedAt, character.CharacterID,
	)

	return err
}

// UpdateAttributes 更新角色属性
func (d *characterDAO) UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error {
	query := `
		UPDATE character_tab SET
			intelligence = ?, emotional_intelligence = ?, memory = ?,
			imagination = ?, physical_fitness = ?, appearance = ?, updated_at = ?
		WHERE character_id = ? AND is_active = 1
	`

	_, err := d.db.ExecContext(ctx, query,
		attributes.Intelligence, attributes.EmotionalIntelligence, attributes.Memory,
		attributes.Imagination, attributes.PhysicalFitness, attributes.Appearance,
		time.Now().UnixMilli(), characterID,
	)

	return err
}

// Delete 删除角色（软删除）
func (d *characterDAO) Delete(ctx context.Context, characterID string) error {
	query := `
		UPDATE character_tab SET
			is_active = 0, updated_at = ?
		WHERE character_id = ? AND is_active = 1
	`

	_, err := d.db.ExecContext(ctx, query, time.Now().UnixMilli(), characterID)
	return err
}

// CountByUserIDAndCharacterID 根据用户ID和角色ID计数（用于权限验证）
func (d *characterDAO) CountByUserIDAndCharacterID(ctx context.Context, userID uint, characterID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM character_tab
		WHERE user_id = ? AND character_id = ? AND is_active = 1
	`

	var count int
	err := d.db.QueryRowContext(ctx, query, userID, characterID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// InsertGameEvent 插入游戏事件
func (d *characterDAO) InsertGameEvent(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO game_events (
			character_id, age, description, impact, created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	result, err := d.db.ExecContext(ctx, query,
		event.CharacterID,
		event.Age,
		event.Description,
		event.Impact,
		event.CreatedAt,
	)
	if err != nil {
		return err
	}

	eventID, err := result.LastInsertId()
	if err == nil {
		event.EventID = uint64(eventID)
	}

	return nil
}

// SelectGameEventsByCharacterID 根据角色ID查询游戏事件
func (d *characterDAO) SelectGameEventsByCharacterID(ctx context.Context, characterID string) ([]models.Event, error) {
	query := `
		SELECT event_id, character_id, age, description, impact, created_at
		FROM game_events
		WHERE character_id = ?
		ORDER BY age ASC, event_id ASC
	`

	rows, err := d.db.QueryContext(ctx, query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.Event, 0)
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
			&event.EventID,
			&event.CharacterID,
			&event.Age,
			&event.Description,
			&event.Impact,
			&event.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

// UpsertGameDecision 新增或更新待处理决策
func (d *characterDAO) UpsertGameDecision(ctx context.Context, decision *models.Decision) error {
	if decision == nil || decision.Options == nil {
		return fmt.Errorf("decision options cannot be nil")
	}

	optionsJSON, err := json.Marshal(decision.Options)
	if err != nil {
		return err
	}

	var previousOptionsJSON []byte
	if decision.PreviousOptions != nil {
		previousOptionsJSON, err = json.Marshal(decision.PreviousOptions)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO game_decision (
			character_id, options, previous_options, previous_decision, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			options = VALUES(options),
			previous_options = VALUES(previous_options),
			previous_decision = VALUES(previous_decision),
			updated_at = VALUES(updated_at)
	`

	_, err = d.db.ExecContext(ctx, query,
		decision.CharacterID,
		optionsJSON,
		previousOptionsJSON,
		decision.PreviousDecision,
		decision.CreatedAt,
		decision.UpdatedAt,
	)
	return err
}

// SelectGameDecisionByCharacterID 查询角色待处理决策
func (d *characterDAO) SelectGameDecisionByCharacterID(ctx context.Context, characterID string) (*models.Decision, error) {
	query := `
		SELECT character_id, options, previous_options, previous_decision, created_at, updated_at
		FROM game_decision
		WHERE character_id = ?
	`

	var decision models.Decision
	var optionsJSON []byte
	var previousOptionsJSON []byte
	var previousDecision sql.NullInt64

	err := d.db.QueryRowContext(ctx, query, characterID).Scan(
		&decision.CharacterID,
		&optionsJSON,
		&previousOptionsJSON,
		&previousDecision,
		&decision.CreatedAt,
		&decision.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var options models.DecisionOption
	if err := json.Unmarshal(optionsJSON, &options); err != nil {
		return nil, err
	}
	decision.Options = &options

	if len(previousOptionsJSON) > 0 {
		var previousOptions models.DecisionOption
		if err := json.Unmarshal(previousOptionsJSON, &previousOptions); err != nil {
			return nil, err
		}
		decision.PreviousOptions = &previousOptions
	}

	if previousDecision.Valid {
		val := int8(previousDecision.Int64)
		decision.PreviousDecision = &val
	}

	return &decision, nil
}

// DeleteGameDecisionByCharacterID 删除角色待处理决策
func (d *characterDAO) DeleteGameDecisionByCharacterID(ctx context.Context, characterID string) error {
	query := `
		DELETE FROM game_decision
		WHERE character_id = ?
	`
	_, err := d.db.ExecContext(ctx, query, characterID)
	return err
}
