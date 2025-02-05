package store

import (
	"context"
	"errors"
	"github.com/HJyup/translatify-common/utils"
	"time"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v5"
)

var ErrMessageNotFound = errors.New("message not found")

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) AddMessage(ctx context.Context, msg *pb.ChatMessage) error {
	query := `
		INSERT INTO chat_messages
			(message_id, from_user_id, to_user_id, content, translated_content, timestamp, translated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := s.dbConn.Exec(ctx, query,
		msg.MessageId,
		msg.FromUserId,
		msg.ToUserId,
		msg.Content,
		"",
		time.Now().Unix(),
		false,
	)
	return err
}

func (s *Store) GetMessage(ctx context.Context, id string) (*pb.ChatMessage, error) {
	query := `
		SELECT message_id, from_user_id, to_user_id, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE message_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)
	return scanChatMessage(row)
}

func (s *Store) ListMessages(ctx context.Context, userID, correspondentID string, since *timestamp.Timestamp) ([]*pb.ChatMessage, error) {
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
	rows, err := s.dbConn.Query(ctx, query, userID, correspondentID, sinceInt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*pb.ChatMessage, 0)
	for rows.Next() {
		msg, err := scanChatMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func scanChatMessage(rs utils.RowScanner) (*pb.ChatMessage, error) {
	var (
		messageID         string
		fromUserID        string
		toUserID          string
		content           string
		translatedContent string
		ts                int64
		translated        bool
	)

	if err := rs.Scan(&messageID, &fromUserID, &toUserID, &content, &translatedContent, &ts, &translated); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &pb.ChatMessage{
		MessageId:         messageID,
		FromUserId:        fromUserID,
		ToUserId:          toUserID,
		Content:           content,
		TranslatedContent: translatedContent,
		Timestamp:         ts,
		Translated:        translated,
	}, nil
}
