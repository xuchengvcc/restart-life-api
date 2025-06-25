package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"github.com/xuchengvcc/restart-life-api/internal/services/ai_services_impl"
)

type AIService interface {
	GenerateText(ctx context.Context, prompt string) (string, error)
}

func NewAIServices(aiConfig config.AIConfig, logger *logrus.Logger) map[string]AIService {
	servicesMap := make(map[string]AIService)
	geminiService, err := ai_services_impl.NewGeminiService(aiConfig.ProviderMap[constants.PROVIDER_GEMINI], logger)
	if err == nil {
		servicesMap[constants.PROVIDER_GEMINI] = geminiService
	}
	tencentService, err := ai_services_impl.NewHunyuanService(aiConfig.ProviderMap[constants.PROVIDER_HUNYUAN], logger)
	if err == nil {
		servicesMap[constants.PROVIDER_HUNYUAN] = tencentService
	}
	return servicesMap
}
