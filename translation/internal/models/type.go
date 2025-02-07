package models

import (
	pb "github.com/HJyup/translatify-common/api"
)

type ConsumerResponse struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	MessageID  string `json:"messageID"`
	Content    string `json:"content"`
}

type TranslationService interface {
	TranslateMessage(sourceLang, targetLang, messageID, content string) (pb.TranslationResponse, error)
}
