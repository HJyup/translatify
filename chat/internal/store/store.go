package store

import (
	"context"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) AddMessage(msg *pb.ChatMessage) error {
	return nil
}

func (s *Store) GetMessage(id string) (*pb.ChatMessage, error) {
	return nil, nil
}

func (s *Store) ListMessages(userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error) {
	messages := make([]*pb.ChatMessage, 0)

	var sinceInt int64 = 0
	if since != nil {
		sinceInt = since.Seconds
	}

	query := `
		SELECT message_id, from_user_id, to_user_id, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE ((from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1))
		  AND timestamp > $3
		ORDER BY timestamp ASC
	`

	rows, err := s.dbConn.Query(context.Background(), query, userID, correspondentID, sinceInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			messageID         string
			fromUserID        string
			toUserID          string
			content           string
			translatedContent string
			ts                int64
			translated        bool
		)

		err = rows.Scan(&messageID, &fromUserID, &toUserID, &content, &translatedContent, &ts, &translated)
		if err != nil {
			return nil, err
		}

		chatMsg := &pb.ChatMessage{
			MessageId:         messageID,
			FromUserId:        fromUserID,
			ToUserId:          toUserID,
			Content:           content,
			TranslatedContent: translatedContent,
			Timestamp:         ts,
			Translated:        translated,
		}
		messages = append(messages, chatMsg)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return messages, nil
}
