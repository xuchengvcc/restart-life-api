package models

type GenerateRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Provider string `json:"provider" binding:"required"`
}
