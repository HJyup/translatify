package models

type CreateChatRequest struct {
	UserAId        string `json:"userAId"`
	UserBId        string `json:"userBId"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}

type SendMessageRequest struct {
	FromUserID string `json:"fromUserId"`
	ToUserID   string `json:"toUserId"`
	Content    string `json:"content"`
}

type ListMessagesRequest struct {
	SinceTimestamp int64  `json:"sinceTimestamp"`
	Limit          int32  `json:"limit"`
	PageToken      string `json:"pageToken"`
}

type StreamMessagesRequest struct {
	ChatId         string `json:"chatId"`
	SinceTimestamp int64  `json:"sinceTimestamp"`
}
