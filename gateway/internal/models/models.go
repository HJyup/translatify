package models

type SendMessageRequest struct {
	FromUserID string `json:"fromUserId"`
	ToUserID   string `json:"toUserId"`
	Content    string `json:"content"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

type ListMessagesRequest struct {
	UserID              string `json:"userId"`
	CorrespondentUserID string `json:"correspondentUserId"`
	SinceTimestamp      int64  `json:"sinceTimestamp"`
}

type StreamMessagesRequest struct {
	UserID              string `json:"userId"`
	CorrespondentUserID string `json:"correspondentUserId"`
	SinceTimestamp      int64  `json:"sinceTimestamp"`
}
