package models

type ConsumerResponse struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	MessageID  string `json:"messageID"`
	Content    string `json:"content"`
}

type TranslationResponse struct {
	TranslatedContent string `json:"translatedContent"`
}

type TranslationService interface {
	TranslateMessage(content, sourceLanguage, targetLanguage string) (*TranslationResponse, error)
}

type TranslatorModel interface {
	TranslateText(text, sourceLanguage, targetLanguage string) (string, error)
}
