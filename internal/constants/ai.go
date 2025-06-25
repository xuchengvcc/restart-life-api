package constants

import "errors"

const (
	PROVIDER_GEMINI  = "gemini"
	PROVIDER_HUNYUAN = "hunyuan"

	DEFAULT_HUNYUAN_MODEL_ID = "hunyuan-lite"
	DEFAULT_GEMINI_MODEL_ID  = "gemini-2.0-flash-lite"
)

var SupportProvider = map[string]struct{}{
	PROVIDER_GEMINI:  {},
	PROVIDER_HUNYUAN: {},
}

var (
	AISERVICE_NOT_INIT     = errors.New("AI service is not initialized")
	AISERVICE_APIKEY_EMPTY = errors.New("AI service API key is empty")
)
