package ai_handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

type AIHandler struct {
	aiServices map[string]services.AIService
	logger     *logrus.Logger
}

func NewAIHandler(aiServices map[string]services.AIService, logger *logrus.Logger) *AIHandler {
	return &AIHandler{
		aiServices: aiServices,
		logger:     logger,
	}
}
