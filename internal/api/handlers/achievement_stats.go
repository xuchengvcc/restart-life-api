package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/api/middleware"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

type AchievementHandler struct {
	service services.AchievementStatsService
	logger  *logrus.Logger
}

type StatsHandler struct {
	service services.AchievementStatsService
	logger  *logrus.Logger
}

func NewAchievementHandler(service services.AchievementStatsService, logger *logrus.Logger) *AchievementHandler {
	return &AchievementHandler{service: service, logger: logger}
}

func NewStatsHandler(service services.AchievementStatsService, logger *logrus.Logger) *StatsHandler {
	return &StatsHandler{service: service, logger: logger}
}

func (h *AchievementHandler) GetCharacterAchievements(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	characterID := c.Param("character_id")
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

	resp, err := h.service.GetAchievements(c.Request.Context(), characterID, userID)
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Warn("failed to get character achievements")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Failed to get achievements",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *AchievementHandler) GetAchievementCategories(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	resp, err := h.service.GetAchievementCategories(c.Request.Context(), userID)
	if err != nil {
		h.logger.WithError(err).Warn("failed to get achievement categories")
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInternalError,
				Message: "Failed to get achievement categories",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *StatsHandler) GetCharacterStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	characterID := c.Param("character_id")
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

	resp, err := h.service.GetCharacterStats(c.Request.Context(), characterID, userID)
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Warn("failed to get character stats")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Failed to get character stats",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
	})
}

func (h *StatsHandler) GetCharacterTimeline(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodePermissionDenied,
				Message: "Unauthorized",
			},
		})
		return
	}

	characterID := c.Param("character_id")
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

	resp, err := h.service.GetCharacterTimeline(c.Request.Context(), characterID, userID)
	if err != nil {
		h.logger.WithError(err).WithField("character_id", characterID).Warn("failed to get character timeline")
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    constants.ErrCodeInvalidParameter,
				Message: "Failed to get character timeline",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    resp,
	})
}
