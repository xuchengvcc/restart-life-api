package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/services"
)

type Handlers struct {
	AuthHandler   *AuthHandler
	HealthHandler *HealthHandler
	AIHandler     *AIHandler
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) InitAuthHandlers(service services.AuthService, logger *logrus.Logger) {
	h.AuthHandler = NewAuthHandler(service, logger)
}

func (h *Handlers) InitHealthHandlers(version string) {
	h.HealthHandler = NewHealthHandler(version)
}

func (h *Handlers) InitAIHandlers(aiServices map[string]services.AIService, logger *logrus.Logger) {
	h.AIHandler = NewAIHandler(aiServices, logger)
}

func (h *Handlers) GetAIHandler() *AIHandler {
	if h == nil {
		return nil
	}
	return h.AIHandler
}

func (h *Handlers) GetAuthHandler() *AuthHandler {
	if h == nil {
		return nil
	}
	return h.AuthHandler
}

func (h *Handlers) GetHealthHandler() *HealthHandler {
	if h == nil {
		return nil
	}
	return h.HealthHandler
}
