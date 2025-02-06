package models

type CreateConversationRequest struct {
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
	ConversationId string `json:"conversationId"`
	SinceTimestamp int64  `json:"sinceTimestamp"`
}
