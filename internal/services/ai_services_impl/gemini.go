package ai_services_impl

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
	"google.golang.org/genai"
)

type GeminiService struct {
	modelIDs []string
	client   *genai.Client
	logger   *logrus.Logger
}

func NewGeminiService(provider *config.Provider, logger *logrus.Logger) (*GeminiService, error) {
	apiKey := provider.GetAPIKey()
	if apiKey == "" {
		logger.WithError(constants.AISERVICE_APIKEY_EMPTY).Error("Failed to initialize Gemini AI service")
		return nil, constants.AISERVICE_APIKEY_EMPTY
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create Gemini AI client")
		return nil, err
	}
	geminiService := &GeminiService{
		client: client,
		logger: logger,
	}
	geminiService.initModelIDs(provider)
	return geminiService, nil
}

func (s *GeminiService) initModelIDs(provider *config.Provider) {
	s.modelIDs = provider.GetModelIDs()
}

func (s *GeminiService) decideModelID() string {
	if len(s.modelIDs) > 0 {
		return s.modelIDs[0]
	}
	return constants.DEFAULT_GEMINI_MODEL_ID
}

func (s *GeminiService) GenerateText(ctx context.Context, prompt string) (string, error) {
	if s == nil || s.client == nil {
		s.logger.WithError(constants.AISERVICE_NOT_INIT).Error("Gemini AI service not initialized")
		return "", constants.AISERVICE_NOT_INIT
	}
	result, err := s.client.Models.GenerateContent(
		ctx,
		s.decideModelID(),
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
