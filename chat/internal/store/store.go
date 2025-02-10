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
var ErrConversationNotFound = errors.New("conversation not found")

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) CreateConversion(ctx context.Context, conv *pb.Conversation) (string, error) {
	query := `
		INSERT INTO conversations
			(conversation_id, username_a, username_b, created_at, source_language, target_language)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ts := time.Now().Unix()
	_, err := s.dbConn.Exec(ctx, query,
		conv.ConversationId,
		conv.UsernameA,
		conv.UsernameB,
		ts,
		conv.SourceLanguage,
		conv.TargetLanguage,
	)
	return conv.ConversationId, err
}

func (s *Store) AddMessage(ctx context.Context, msg *pb.ChatMessage) error {
	query := `
		INSERT INTO chat_messages
			(message_id, conversation_id, sender_username, receiver_username, content, translated_content, timestamp, translated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	ts := time.Now().Unix()
	_, err := s.dbConn.Exec(ctx, query,
		msg.MessageId,
		msg.ConversationId,
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
		SELECT message_id, conversation_id, sender_username, receiver_username, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE message_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)
	return scanChatMessage(row)
}

func (s *Store) ListMessages(ctx context.Context, conversationID string, since *timestamp.Timestamp, limit int, pageToken string) ([]*pb.ChatMessage, string, error) {
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
		SELECT message_id, conversation_id, sender_username, receiver_username, content, translated_content, timestamp, translated
		FROM chat_messages
		WHERE conversation_id = $1 AND timestamp > $2
		ORDER BY timestamp ASC
		LIMIT $3
	`
	rows, err := s.dbConn.Query(ctx, query, conversationID, effectiveSince, limit+1)
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

func (s *Store) GetConversation(ctx context.Context, id string) (*pb.Conversation, error) {
	query := `
		SELECT conversation_id, username_a, username_b, created_at, source_language, target_language
		FROM conversations
		WHERE conversation_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)

	var (
		conversationID string
		userNameA      string
		userNameB      string
		createdAt      int64
		sourceLang     string
		targetLang     string
	)
	if err := row.Scan(&conversationID, &userNameA, &userNameB, &createdAt, &sourceLang, &targetLang); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	return &pb.Conversation{
		ConversationId: conversationID,
		UsernameA:      userNameA,
		UsernameB:      userNameB,
		CreatedAt:      createdAt,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}, nil
}

func (s *Store) ListConversations(ctx context.Context, userID string) ([]*pb.Conversation, error) {
	query := `
		SELECT conversation_id, username_a, username_b, created_at, source_language, target_language
		FROM conversations
		WHERE user_a_id = $1 OR user_b_id = $1
	`
	rows, err := s.dbConn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := make([]*pb.Conversation, 0)
	for rows.Next() {
		var (
			conversationID string
			userNameA      string
			userNameB      string
			createdAt      int64
			sourceLang     string
			targetLang     string
		)
		if err = rows.Scan(&conversationID, &userNameA, &userNameB, &createdAt, &sourceLang, &targetLang); err != nil {
			return nil, err
		}
		conversations = append(conversations, &pb.Conversation{
			ConversationId: conversationID,
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

	return conversations, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanChatMessage(rs rowScanner) (*pb.ChatMessage, error) {
	var (
		messageID         string
		conversationID    string
		senderUserName    string
		receiverUserName  string
		content           string
		translatedContent string
		ts                int64
		translated        bool
	)

	if err := rs.Scan(&messageID, &conversationID, &senderUserName, &receiverUserName, &content, &translatedContent, &ts, &translated); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &pb.ChatMessage{
		MessageId:         messageID,
		ConversationId:    conversationID,
		SenderUsername:    senderUserName,
		ReceiverUsername:  receiverUserName,
		Content:           content,
		TranslatedContent: translatedContent,
		Timestamp:         ts,
		Translated:        translated,
	}, nil
}
