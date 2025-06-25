package ai_services_impl

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
	"github.com/xuchengvcc/restart-life-api/internal/config"
	"github.com/xuchengvcc/restart-life-api/internal/constants"
)

type HunyuanService struct {
	modelIDs []string // 可能需要的模型ID列表
	client   *openai.Client
	logger   *logrus.Logger
}

func NewHunyuanService(provider *config.Provider, logger *logrus.Logger) (*HunyuanService, error) {
	apiKey := provider.GetAPIKey()
	if apiKey == "" {
		logger.WithError(constants.AISERVICE_APIKEY_EMPTY).Error("Failed to initialize Tencent AI service")
		return nil, constants.AISERVICE_APIKEY_EMPTY
	}
	client := openai.NewClient(
		option.WithAPIKey(apiKey), // 混元 APIKey
		option.WithBaseURL("https://api.hunyuan.cloud.tencent.com/v1/"), // 混元 endpoint
	)
	hunyuanService := &HunyuanService{
		client: &client,
		logger: logger,
	}
	hunyuanService.initModelIDs(provider)
	return hunyuanService, nil
}

func (s *HunyuanService) initModelIDs(provider *config.Provider) {
	s.modelIDs = provider.GetModelIDs()
}

func (s *HunyuanService) decideModelID() string {
	if len(s.modelIDs) > 0 {
		return s.modelIDs[0] // 返回第一个模型ID，或根据其他逻辑选择
	}
	return constants.DEFAULT_HUNYUAN_MODEL_ID
}

func (s *HunyuanService) GenerateText(ctx context.Context, prompt string) (string, error) {
	if s == nil || s.client == nil {
		s.logger.WithError(constants.AISERVICE_NOT_INIT).Error("Tencent AI service not initialized")
		return "", constants.AISERVICE_NOT_INIT
	}
	chatCompletion, err := s.client.Chat.Completions.New(ctx,
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
			Model: s.decideModelID(),
		},
		option.WithJSONSet("enable_enhancement", true), // <- 自定义参数
	)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate text with Tencent AI service")
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
