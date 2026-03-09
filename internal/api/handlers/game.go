package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

type GameHandler struct {
	gameService services.GameService
	logger      *logrus.Logger
}

func NewGameHandler(gameService services.GameService, logger *logrus.Logger) *GameHandler {
	return &GameHandler{
		gameService: gameService,
		logger:      logger,
	}
}

// StartOrResumeGame 开始或继续游戏
// @Summary 开始或继续游戏
// @Description 优先返回未完成角色的游戏状态，否则引导新建角色
// @Tags Game
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=models.GameState}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/start [post]
func (h *GameHandler) StartOrResumeGame(c *gin.Context) {
	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "用户ID格式错误",
			},
		})
		return
	}

	// 将字符串转换为 uint64
	userIDuint64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "用户ID格式错误",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": userIDuint64,
	}).Info("开始或继续游戏请求")

	// 开始或继续游戏
	gameState, err := h.gameService.StartOrResumeGame(c.Request.Context(), userIDuint64)
	if err != nil {
		h.logger.WithError(err).Error("开始或继续游戏失败")
		// 如果错误是需要创建新角色，返回特定错误码
		if err.Error() == "需要创建新角色" {
			c.JSON(http.StatusOK, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    constants.ErrCodeNeedCreateCharacter,
					Message: "需要创建新角色",
					Details: "没有找到未完成的角色，请先创建角色",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "开始或继续游戏失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameState,
	})
}

// StartGame 为指定角色开始游戏（兼容旧API）
// @Summary 为指定角色开始游戏
// @Description 为指定角色开始一局新游戏（已废弃，建议使用StartOrResumeGame）
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Success 200 {object} models.APIResponse{data=models.GameState}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/start/{character_id} [post]
func (h *GameHandler) StartGame(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("开始游戏请求")

	// 获取游戏状态
	gameState, err := h.gameService.GetGameState(c.Request.Context(), characterID)
	if err != nil {
		h.logger.WithError(err).Error("获取游戏状态失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "获取游戏状态失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameState,
	})
}

// AdvanceGame 推进游戏（智能处理决策和推进）
// @Summary 推进游戏
// @Description 智能推进游戏：如果有待处理决策则需要option_type，否则自动推进年龄并生成新事件
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Param request body models.GameProgressRequest false "游戏推进请求（决策时需要option_type）"
// @Success 200 {object} models.APIResponse{data=models.GameState}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/advance/{character_id} [post]
func (h *GameHandler) AdvanceGame(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 绑定请求参数（可选）
	var req models.GameProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果没有请求体或解析失败，设为默认值
		req = models.GameProgressRequest{}
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
		"option_type":  req.OptionType,
	}).Info("推进游戏请求")

	// 为推进游戏请求设置超时，避免长时间阻塞
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// 智能推进游戏
	gameState, err := h.gameService.AdvanceGameSmart(ctx, characterID, req.OptionType)
	if err != nil {
		h.logger.WithError(err).Error("推进游戏失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "推进游戏失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameState,
	})
}

// GetGameState 获取游戏状态
// @Summary 获取游戏状态
// @Description 获取角色的当前游戏状态
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Success 200 {object} models.APIResponse{data=models.GameState}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/state/{character_id} [get]
func (h *GameHandler) GetGameState(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("获取游戏状态请求")

	// 获取游戏状态
	gameState, err := h.gameService.GetGameState(c.Request.Context(), characterID)
	if err != nil {
		h.logger.WithError(err).Error("获取游戏状态失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "获取游戏状态失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameState,
	})
}

// SaveGame 保存游戏
// @Summary 保存游戏
// @Description 保存角色的游戏进度
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/save/{character_id} [post]
func (h *GameHandler) SaveGame(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("保存游戏请求")

	// 保存游戏
	err := h.gameService.SaveGame(c.Request.Context(), characterID)
	if err != nil {
		h.logger.WithError(err).Error("保存游戏失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "保存游戏失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    nil,
	})
}

// LoadGame 加载游戏
// @Summary 加载游戏
// @Description 加载角色的游戏进度
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Success 200 {object} models.APIResponse{data=models.GameState}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/load/{character_id} [post]
func (h *GameHandler) LoadGame(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
	}).Info("加载游戏请求")

	// 加载游戏
	gameState, err := h.gameService.LoadGame(c.Request.Context(), characterID)
	if err != nil {
		h.logger.WithError(err).Error("加载游戏失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "加载游戏失败",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameState,
	})
}

// GetEventHistory 获取事件历史
// @Summary 获取事件历史
// @Description 获取角色的人生事件历史记录
// @Tags Game
// @Accept json
// @Produce json
// @Param character_id path string true "角色ID"
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Success 200 {object} models.APIResponse{data=[]models.Event}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/game/events/{character_id} [get]
func (h *GameHandler) GetEventHistory(c *gin.Context) {
	characterID := c.Param("character_id")
	if characterID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "角色ID不能为空",
			},
		})
		return
	}

	// 获取分页参数
	page := 1
	limit := 20
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// 获取用户ID（从JWT中间件获取）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "用户未认证",
			},
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":      userID,
		"character_id": characterID,
		"page":         page,
		"limit":        limit,
	}).Info("获取事件历史请求")

	// 获取事件历史
	events, err := h.gameService.GetEventHistory(c.Request.Context(), characterID)
	if err != nil {
		h.logger.WithError(err).Error("获取事件历史失败")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "获取事件历史失败",
				Details: err.Error(),
			},
		})
		return
	}

	// 简单的分页处理（实际应该在服务层处理）
	start := (page - 1) * limit
	end := start + limit
	if start > len(events) {
		events = []models.Event{}
	} else if end > len(events) {
		events = events[start:]
	} else {
		events = events[start:end]
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    events,
	})
}
