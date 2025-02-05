package models

import (
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type TranslationService interface {
	TranslateMessage(msgID, content, sourceLan, targetLan string, time timestamp.Timestamp) (pb.TranslationResponse, error)
}

type ChatStore interface {
}
