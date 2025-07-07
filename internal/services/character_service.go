package services

import (
	"context"
	"fmt"

	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/repository"
)

// CharacterService 角色服务接口
type CharacterService interface {
	CreateCharacter(ctx context.Context, userID uint, req *models.CreateCharacterRequest) (*models.Character, error)
	GetCharacter(ctx context.Context, characterID string, userID uint) (*models.Character, error)
	GetUserCharacters(ctx context.Context, userID uint) (*models.CharacterListResponse, error)
	GetActiveCharacters(ctx context.Context, userID uint) (*models.CharacterListResponse, error)
	UpdateCharacter(ctx context.Context, characterID string, userID uint, req *models.UpdateCharacterRequest) (*models.Character, error)
	UpdateCharacterAttributes(ctx context.Context, characterID string, userID uint, req *models.UpdateCharacterAttributesRequest) (*models.Character, error)
	DeleteCharacter(ctx context.Context, characterID string, userID uint) error
	ValidateOwnership(ctx context.Context, characterID string, userID uint) error
}

// characterService 角色服务实现
type characterService struct {
	characterRepo repository.CharacterRepository
}

// NewCharacterService 创建角色服务
func NewCharacterService(characterRepo repository.CharacterRepository) CharacterService {
	return &characterService{
		characterRepo: characterRepo,
	}
}

// CreateCharacter 创建角色
func (s *characterService) CreateCharacter(ctx context.Context, userID uint, req *models.CreateCharacterRequest) (*models.Character, error) {
	// 创建角色对象
	character := &models.Character{
		UserID:        userID,
		CharacterName: req.CharacterName,
		BirthCountry:  req.BirthCountry,
		BirthYear:     req.BirthYear,
		Gender:        req.Gender,
		Race:          req.Race,
	}

	// 生成随机属性
	character.Attributes = s.characterRepo.GenerateRandomAttributes()
	character.SetDefaultValues()

	// 保存到数据库
	err := s.characterRepo.Create(ctx, character)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %w", err)
	}

	return character, nil
}

// GetCharacter 获取角色详情
func (s *characterService) GetCharacter(ctx context.Context, characterID string, userID uint) (*models.Character, error) {
	// 验证所有权
	err := s.ValidateOwnership(ctx, characterID, userID)
	if err != nil {
		return nil, err
	}

	// 获取角色
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	return character, nil
}

// GetUserCharacters 获取用户所有角色
func (s *characterService) GetUserCharacters(ctx context.Context, userID uint) (*models.CharacterListResponse, error) {
	characters, err := s.characterRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user characters: %w", err)
	}

	summaries := s.characterRepo.CreateCharacterSummaries(characters)

	return &models.CharacterListResponse{
		Characters: summaries,
		Total:      len(summaries),
	}, nil
}

// GetActiveCharacters 获取用户活跃角色
func (s *characterService) GetActiveCharacters(ctx context.Context, userID uint) (*models.CharacterListResponse, error) {
	characters, err := s.characterRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active characters: %w", err)
	}

	summaries := s.characterRepo.CreateCharacterSummaries(characters)

	return &models.CharacterListResponse{
		Characters: summaries,
		Total:      len(summaries),
	}, nil
}

// UpdateCharacter 更新角色信息
func (s *characterService) UpdateCharacter(ctx context.Context, characterID string, userID uint, req *models.UpdateCharacterRequest) (*models.Character, error) {
	// 验证所有权
	err := s.ValidateOwnership(ctx, characterID, userID)
	if err != nil {
		return nil, err
	}

	// 获取当前角色
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// 更新字段
	s.updateCharacterFields(character, req)

	// 保存更新
	err = s.characterRepo.Update(ctx, character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character: %w", err)
	}

	return character, nil
}

// UpdateCharacterAttributes 更新角色属性
func (s *characterService) UpdateCharacterAttributes(ctx context.Context, characterID string, userID uint, req *models.UpdateCharacterAttributesRequest) (*models.Character, error) {
	// 验证所有权
	err := s.ValidateOwnership(ctx, characterID, userID)
	if err != nil {
		return nil, err
	}

	// 获取当前角色
	character, err := s.characterRepo.GetByID(ctx, characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// 更新属性
	s.updateCharacterAttributes(&character.Attributes, req)

	// 保存属性更新
	err = s.characterRepo.UpdateAttributes(ctx, characterID, &character.Attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to update character attributes: %w", err)
	}

	return character, nil
}

// DeleteCharacter 删除角色
func (s *characterService) DeleteCharacter(ctx context.Context, characterID string, userID uint) error {
	// 验证所有权
	err := s.ValidateOwnership(ctx, characterID, userID)
	if err != nil {
		return err
	}

	// 删除角色
	err = s.characterRepo.Delete(ctx, characterID)
	if err != nil {
		return fmt.Errorf("failed to delete character: %w", err)
	}

	return nil
}

// ValidateOwnership 验证角色所有权
func (s *characterService) ValidateOwnership(ctx context.Context, characterID string, userID uint) error {
	isOwner, err := s.characterRepo.IsOwner(ctx, characterID, userID)
	if err != nil {
		return fmt.Errorf("failed to validate ownership: %w", err)
	}

	if !isOwner {
		return fmt.Errorf("character not found or permission denied")
	}

	return nil
}

// updateCharacterFields 更新角色字段
func (s *characterService) updateCharacterFields(character *models.Character, req *models.UpdateCharacterRequest) {
	if req.CharacterName != nil {
		character.CharacterName = *req.CharacterName
	}
	if req.EducationLevel != nil {
		character.EducationLevel = req.EducationLevel
	}
	if req.MaritalStatus != nil {
		character.MaritalStatus = req.MaritalStatus
	}
	if req.CurrentCountry != nil {
		character.CurrentCountry = req.CurrentCountry
	}
	if req.CurrentLocation != nil {
		character.CurrentLocation = req.CurrentLocation
	}
	if req.CurrentActivity != nil {
		character.CurrentActivity = req.CurrentActivity
	}
	if req.Personality != nil {
		character.Personality = req.Personality
	}
	if req.Career != nil {
		character.Career = req.Career
	}
	if req.SkillTendency != nil {
		character.SkillTendency = req.SkillTendency
	}
	if req.FamilyBackground != nil {
		character.FamilyBackground = req.FamilyBackground
	}
	if req.SocialRelationships != nil {
		character.SocialRelationships = req.SocialRelationships
	}
	if req.CareerDesc != nil {
		character.CareerDesc = req.CareerDesc
	}
	if req.EducationDesc != nil {
		character.EducationDesc = req.EducationDesc
	}
}

// updateCharacterAttributes 更新角色属性
func (s *characterService) updateCharacterAttributes(attributes *models.CharacterAttributes, req *models.UpdateCharacterAttributesRequest) {
	if req.Intelligence != nil {
		attributes.Intelligence = *req.Intelligence
	}
	if req.EmotionalIntelligence != nil {
		attributes.EmotionalIntelligence = *req.EmotionalIntelligence
	}
	if req.Memory != nil {
		attributes.Memory = *req.Memory
	}
	if req.Imagination != nil {
		attributes.Imagination = *req.Imagination
	}
	if req.PhysicalFitness != nil {
		attributes.PhysicalFitness = *req.PhysicalFitness
	}
	if req.Appearance != nil {
		attributes.Appearance = *req.Appearance
	}
}
