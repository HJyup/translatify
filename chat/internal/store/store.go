package store

import (
	"context"
	"errors"
	"strconv"
	"time"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v5"
)

var ErrMessageNotFound = errors.New("message not found")
var ErrChatNotFound = errors.New("chat not found")

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) CreateConversion(ctx context.Context, conv *pb.Chat) (string, error) {
	query := `
		INSERT INTO chats
			(chat_id, username_a, username_b, created_at, source_language, target_language)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ts := time.Now().Unix()
	_, err := s.dbConn.Exec(ctx, query,
		conv.ChatId,
		conv.UsernameA,
		conv.UsernameB,
		ts,
		conv.SourceLanguage,
		conv.TargetLanguage,
	)
	return conv.ChatId, err
}

func (s *Store) AddMessage(ctx context.Context, msg *pb.ChatMessage) error {
	query := `
		INSERT INTO chat_messages
			(message_id, chat_id, sender_username, receiver_username, content, translated_content, timestamp, translated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	ts := time.Now().Unix()
	_, err := s.dbConn.Exec(ctx, query,
		msg.MessageId,
		msg.ChatId,
		msg.SenderUsername,
		msg.ReceiverUsername,
		msg.Content,
		"",
		ts,
		false,
	)
	return err
}

func (s *Store) GetMessage(ctx context.Context, id string) (*pb.ChatMessage, error) {
	query := `
		SELECT message_id, chat_id, sender_username, receiver_username, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE message_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)
	return scanChatMessage(row)
}

func (s *Store) ListMessages(ctx context.Context, chatId string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error) {
	var effectiveSince int64 = 0
	if since != nil {
		effectiveSince = since.Seconds
	}
	if pageToken != "" {
		if tokenTs, err := strconv.ParseInt(pageToken, 10, 64); err == nil && tokenTs > effectiveSince {
			effectiveSince = tokenTs
		}
	}

	query := `
		SELECT message_id, chat_id, sender_username, receiver_username, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE chat_id = $1 AND timestamp > $2
		ORDER BY timestamp ASC
		LIMIT $3
	`
	rows, err := s.dbConn.Query(ctx, query, chatId, effectiveSince, limit+1)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	messages := make([]*pb.ChatMessage, 0)
	for rows.Next() {
		msg, err := scanChatMessage(rows)
		if err != nil {
			return nil, "", err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(messages) > limit {
		nextPageToken = strconv.FormatInt(messages[limit].Timestamp, 10)
		messages = messages[:limit]
	}

	return messages, nextPageToken, nil
}

func (s *Store) GetChat(ctx context.Context, id string) (*pb.Chat, error) {
	query := `
		SELECT chat_id, username_a, username_b, created_at, source_language, target_language
		FROM chats
		WHERE chat_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)

	var (
		ChatID     string
		userNameA  string
		userNameB  string
		createdAt  int64
		sourceLang string
		targetLang string
	)
	if err := row.Scan(&ChatID, &userNameA, &userNameB, &createdAt, &sourceLang, &targetLang); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}

	return &pb.Chat{
		ChatId:         ChatID,
		UsernameA:      userNameA,
		UsernameB:      userNameB,
		CreatedAt:      createdAt,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}, nil
}

func (s *Store) ListChats(ctx context.Context, userID string) ([]*pb.Chat, error) {
	query := `
		SELECT chat_id, username_a, username_b, created_at, source_language, target_language
		FROM chats
		WHERE user_a_id = $1 OR user_b_id = $1
	`
	rows, err := s.dbConn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Chats := make([]*pb.Chat, 0)
	for rows.Next() {
		var (
			chatID     string
			userNameA  string
			userNameB  string
			createdAt  int64
			sourceLang string
			targetLang string
		)
		if err = rows.Scan(&chatID, &userNameA, &userNameB, &createdAt, &sourceLang, &targetLang); err != nil {
			return nil, err
		}
		Chats = append(Chats, &pb.Chat{
			ChatId:         chatID,
			UsernameA:      userNameA,
			UsernameB:      userNameB,
			CreatedAt:      createdAt,
			SourceLanguage: sourceLang,
			TargetLanguage: targetLang,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Chats, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanChatMessage(rs rowScanner) (*pb.ChatMessage, error) {
	var (
		messageID         string
		chatID            string
		senderUserName    string
		receiverUserName  string
		content           string
		translatedContent string
		ts                int64
		translated        bool
	)

	if err := rs.Scan(&messageID, &chatID, &senderUserName, &receiverUserName, &content, &translatedContent, &ts, &translated); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &pb.ChatMessage{
		MessageId:         messageID,
		ChatId:            chatID,
		SenderUsername:    senderUserName,
		ReceiverUsername:  receiverUserName,
		Content:           content,
		TranslatedContent: translatedContent,
		Timestamp:         ts,
		Translated:        translated,
	}, nil
}
