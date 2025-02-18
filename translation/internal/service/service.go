package service

import (
	"fmt"
	"github.com/HJyup/translatify-translation/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TranslationService struct {
	translator models.TranslatorModel

	channel *amqp.Channel
}

func NewTranslationService(translator models.TranslatorModel) *TranslationService {
	return &TranslationService{translator: translator}
}

func (s *TranslationService) TranslateMessage(sourceLang, targetLang, content string) (*models.TranslationResponse, error) {
	translatedText, err := s.translator.TranslateText(content, sourceLang, targetLang)
	if err != nil {
		return nil, fmt.Errorf("translation error: %w", err)
	}

	response := &models.TranslationResponse{
		TranslatedContent: translatedText,
	}

	return response, nil
}
