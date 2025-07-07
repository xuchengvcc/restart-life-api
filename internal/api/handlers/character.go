package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

// CharacterHandler 角色处理器
type CharacterHandler struct {
	characterService services.CharacterService
	logger           *logrus.Logger
}

// NewCharacterHandler 创建角色处理器
func NewCharacterHandler(characterService services.CharacterService, logger *logrus.Logger) *CharacterHandler {
	return &CharacterHandler{
		characterService: characterService,
		logger:           logger,
	}
}

// CreateCharacter 创建角色
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req models.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for create character")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 创建角色
	character, err := h.characterService.CreateCharacter(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create character")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to create character",
				Details: err.Error(),
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": character.CharacterID,
	}).Info("Character created successfully")

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    character,
	})
}

// GetCharacter 获取角色详情
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Character ID is required",
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 获取角色
	character, err := h.characterService.GetCharacter(c.Request.Context(), characterID, userID.(uint))
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Error("Failed to get character")
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeResourceNotFound,
				Message: "Character not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    character,
	})
}

// GetUserCharacters 获取用户角色列表
func (h *CharacterHandler) GetUserCharacters(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 检查是否只获取活跃角色
	activeOnly := c.Query("active") == "true"

	var response *models.CharacterListResponse
	var err error

	if activeOnly {
		response, err = h.characterService.GetActiveCharacters(c.Request.Context(), userID.(uint))
	} else {
		response, err = h.characterService.GetUserCharacters(c.Request.Context(), userID.(uint))
	}

	if err != nil {
		h.logger.WithError(err).Error("Failed to get user characters")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to get characters",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    response,
	})
}

// UpdateCharacter 更新角色信息
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Character ID is required",
			},
		})
		return
	}

	var req models.UpdateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for update character")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 更新角色
	character, err := h.characterService.UpdateCharacter(c.Request.Context(), characterID, userID.(uint), &req)
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Error("Failed to update character")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to update character",
				Details: err.Error(),
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("Character updated successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    character,
	})
}

// UpdateCharacterAttributes 更新角色属性
func (h *CharacterHandler) UpdateCharacterAttributes(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Character ID is required",
			},
		})
		return
	}

	var req models.UpdateCharacterAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid request body for update character attributes")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Invalid request body",
				Details: err.Error(),
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 更新角色属性
	character, err := h.characterService.UpdateCharacterAttributes(c.Request.Context(), characterID, userID.(uint), &req)
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Error("Failed to update character attributes")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to update character attributes",
				Details: err.Error(),
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("Character attributes updated successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    character,
	})
}

// DeleteCharacter 删除角色
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Character ID is required",
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 删除角色
	err := h.characterService.DeleteCharacter(c.Request.Context(), characterID, userID.(uint))
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Error("Failed to delete character")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to delete character",
				Details: err.Error(),
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("Character deleted successfully")

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gin.H{"message": "Character deleted successfully"},
	})
}

// GetCharacterAttributes 获取角色属性
func (h *CharacterHandler) GetCharacterAttributes(c *gin.Context) {
	characterID := c.Param("id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Character ID is required",
			},
		})
		return
	}

	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	// 获取角色
	character, err := h.characterService.GetCharacter(c.Request.Context(), characterID, userID.(uint))
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Error("Failed to get character")
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeResourceNotFound,
				Message: "Character not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    character.Attributes,
	})
}
