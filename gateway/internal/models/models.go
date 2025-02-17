package models

type CreateChatRequest struct {
	UserNameA      string `json:"usernameA"`
	UserNameB      string `json:"userNameB"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage"`
}

type SendMessageRequest struct {
	FromUserName string `json:"fromUsername"`
	ToUserName   string `json:"toUsername"`
	Content      string `json:"content"`
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

type CreateUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Language string `json:"language"`
}
