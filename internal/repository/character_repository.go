package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/xuchengvcc/restart-life-api/internal/dao"
	"github.com/xuchengvcc/restart-life-api/internal/models"
)

// CharacterRepository 角色仓储接口
type CharacterRepository interface {
	Create(ctx context.Context, character *models.Character) error
	GetByID(ctx context.Context, characterID string) (*models.Character, error)
	GetByUserID(ctx context.Context, userID uint) ([]*models.Character, error)
	GetActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error
	Delete(ctx context.Context, characterID string) error
	IsOwner(ctx context.Context, characterID string, userID uint) (bool, error)
	GenerateRandomAttributes() models.CharacterAttributes
	CreateCharacterSummaries(characters []*models.Character) []models.CharacterSummary
}

// characterRepository 角色仓储实现
type characterRepository struct {
	characterDAO dao.CharacterDAO
}

// NewCharacterRepository 创建角色仓储
func NewCharacterRepository(characterDAO dao.CharacterDAO) CharacterRepository {
	return &characterRepository{
		characterDAO: characterDAO,
	}
}

// Create 创建角色
func (r *characterRepository) Create(ctx context.Context, character *models.Character) error {
	// 生成角色ID (简单的UUID替代)
	character.CharacterID = r.generateCharacterID()

	return r.characterDAO.Insert(ctx, character)
}

// GetByID 根据ID获取角色
func (r *characterRepository) GetByID(ctx context.Context, characterID string) (*models.Character, error) {
	return r.characterDAO.SelectByID(ctx, characterID)
}

// GetByUserID 根据用户ID获取角色列表
func (r *characterRepository) GetByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	return r.characterDAO.SelectByUserID(ctx, userID)
}

// GetActiveByUserID 根据用户ID获取活跃角色列表
func (r *characterRepository) GetActiveByUserID(ctx context.Context, userID uint) ([]*models.Character, error) {
	return r.characterDAO.SelectActiveByUserID(ctx, userID)
}

// Update 更新角色信息
func (r *characterRepository) Update(ctx context.Context, character *models.Character) error {
	character.UpdatedAt = time.Now().UnixMilli()
	return r.characterDAO.Update(ctx, character)
}

// UpdateAttributes 更新角色属性
func (r *characterRepository) UpdateAttributes(ctx context.Context, characterID string, attributes *models.CharacterAttributes) error {
	return r.characterDAO.UpdateAttributes(ctx, characterID, attributes)
}

// Delete 删除角色（软删除）
func (r *characterRepository) Delete(ctx context.Context, characterID string) error {
	return r.characterDAO.Delete(ctx, characterID)
}

// IsOwner 检查角色是否属于指定用户
func (r *characterRepository) IsOwner(ctx context.Context, characterID string, userID uint) (bool, error) {
	count, err := r.characterDAO.CountByUserIDAndCharacterID(ctx, userID, characterID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GenerateRandomAttributes 生成随机属性
func (r *characterRepository) GenerateRandomAttributes() models.CharacterAttributes {
	// 属性范围：30-70，平均50
	minAttr := 30
	maxAttr := 70

	return models.CharacterAttributes{
		Intelligence:          r.randomAttribute(minAttr, maxAttr),
		EmotionalIntelligence: r.randomAttribute(minAttr, maxAttr),
		Memory:                r.randomAttribute(minAttr, maxAttr),
		Imagination:           r.randomAttribute(minAttr, maxAttr),
		PhysicalFitness:       r.randomAttribute(minAttr, maxAttr),
		Appearance:            r.randomAttribute(minAttr, maxAttr),
	}
}

// CreateCharacterSummaries 创建角色概要列表
func (r *characterRepository) CreateCharacterSummaries(characters []*models.Character) []models.CharacterSummary {
	summaries := make([]models.CharacterSummary, len(characters))

	for i, character := range characters {
		summaries[i] = models.CharacterSummary{
			CharacterID:   character.CharacterID,
			CharacterName: character.CharacterName,
			CurrentAge:    character.CurrentAge,
			Gender:        character.Gender,
			Race:          character.Race,
			LifeStage:     character.LifeStage,
			Attributes:    character.Attributes,
			CreatedAt:     character.CreatedAt,
			IsActive:      character.IsActive,
			Summary:       character.Summary,
		}
	}

	return summaries
}

// generateCharacterID 生成角色ID
func (r *characterRepository) generateCharacterID() string {
	// 简单的ID生成：时间戳 + 随机数
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("char_%d_%d", timestamp, randomNum)
}

// randomAttribute 生成随机属性值
func (r *characterRepository) randomAttribute(min, max int) int {
	return rand.Intn(max-min+1) + min
}
