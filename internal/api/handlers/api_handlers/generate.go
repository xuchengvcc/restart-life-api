package ai_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/models"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

func (h *AIHandler) Generate(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid generate request")
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if errMsg := h.CheckGenerateRequest(&req); errMsg != "" {
		h.logger.Warnf("Generate request param invalid: %s", errMsg)
		response := models.NewErrorResponse(models.ErrCodeValidationFailed, "Invalid request data", errMsg)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	res, err := h.getService(req.Provider).GenerateText(c.Request.Context(), req.Prompt)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate text")
		response := models.NewErrorResponse(models.ErrCodeInternalError, "Failed to generate text", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.NewSuccessResponse(res)
	c.JSON(http.StatusOK, response)
}

func (h *AIHandler) getService(provider string) services.AIService {
	if h.aiServices == nil {
		h.logger.Error("AI services not initialized")
		return nil
	}
	return h.aiServices[provider]
}

func (h *AIHandler) CheckGenerateRequest(req *models.GenerateRequest) string {
	if req.Prompt == "" {
		return "prompt is empty"
	}

	if _, exist := constants.SupportProvider[req.Provider]; !exist {
		return "provider not supported"
	}
	return ""
}
