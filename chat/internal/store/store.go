package store

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"
	"time"

	"github.com/HJyup/translatify-chat/internal/models"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) CreateConversion(ctx context.Context, conv *models.Chat) (string, error) {
	query := `
		INSERT INTO chats
			(username_a, username_b, created_at, source_language, target_language)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING chat_id
	`
	now := time.Now()
	conv.CreatedAt = now
	var chatID string
	err := s.dbConn.QueryRow(ctx, query,
		conv.UsernameA,
		conv.UsernameB,
		now.Unix(),
		conv.SourceLang,
		conv.TargetLang,
	).Scan(&chatID)
	if err != nil {
		return "", err
	}
	return chatID, nil
}

func (s *Store) AddMessage(ctx context.Context, msg *models.ChatMessage) (string, error) {
	ctx, span := otel.Tracer("chat-store").Start(ctx, "AddMessage")
	span.SetAttributes(attribute.String("chatID", msg.ChatID))
	defer span.End()

	query := `
		INSERT INTO messages
			(chat_id, sender_username, receiver_username, content, translated_content, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING message_id
	`
	now := time.Now()
	msg.Timestamp = now
	var messageID string
	err := s.dbConn.QueryRow(ctx, query,
		msg.ChatID,
		msg.SenderUsername,
		msg.ReceiverUsername,
		msg.Content,
		"",
		now.Unix(),
	).Scan(&messageID)
	if err != nil {
		return "", err
	}
	return messageID, nil
}

func (s *Store) GetMessage(ctx context.Context, id string) (*models.ChatMessage, error) {
	query := `
		SELECT message_id, chat_id, sender_username, receiver_username, content, translated_content, timestamp
		FROM messages
		WHERE message_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)
	return scanChatMessage(row)
}

func (s *Store) ListMessages(ctx context.Context, chatID string, since *time.Time, limit int, pageToken string) ([]*models.ChatMessage, string, error) {
	effectiveSince := int64(0)
	if since != nil {
		effectiveSince = since.Unix()
	}
	if pageToken != "" {
		if tokenTs, err := strconv.ParseInt(pageToken, 10, 64); err == nil && tokenTs > effectiveSince {
			effectiveSince = tokenTs
		}
	}

	query := `
		SELECT message_id, chat_id, sender_username, receiver_username, content, translated_content, timestamp
		FROM messages
		WHERE chat_id = $1 AND timestamp > $2
		ORDER BY timestamp ASC
		LIMIT $3
	`
	rows, err := s.dbConn.Query(ctx, query, chatID, effectiveSince, limit+1)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	messages := make([]*models.ChatMessage, 0)
	for rows.Next() {
		msg, err := scanChatMessage(rows)
		if err != nil {
			return nil, "", err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(messages) > limit {
		nextPageToken = strconv.FormatInt(messages[limit].Timestamp.Unix(), 10)
		messages = messages[:limit]
	}

	return messages, nextPageToken, nil
}

func (s *Store) GetChat(ctx context.Context, id string) (*models.Chat, error) {
	query := `
		SELECT chat_id, username_a, username_b, created_at, source_language, target_language
		FROM chats
		WHERE chat_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, id)

	var (
		chatID     string
		usernameA  string
		usernameB  string
		createdAt  int64
		sourceLang string
		targetLang string
	)
	if err := row.Scan(&chatID, &usernameA, &usernameB, &createdAt, &sourceLang, &targetLang); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("chat not found")
		}
		return nil, err
	}

	return &models.Chat{
		ChatID:     chatID,
		UsernameA:  usernameA,
		UsernameB:  usernameB,
		CreatedAt:  time.Unix(createdAt, 0),
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}, nil
}

func (s *Store) ListChats(ctx context.Context, userName string) ([]*models.Chat, error) {
	query := `
		SELECT chat_id, username_a, username_b, created_at, source_language, target_language
		FROM chats
		WHERE username_a = $1 OR username_b = $1
	`
	rows, err := s.dbConn.Query(ctx, query, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*models.Chat, 0)
	for rows.Next() {
		var (
			chatID     string
			usernameA  string
			usernameB  string
			createdAt  int64
			sourceLang string
			targetLang string
		)
		if err = rows.Scan(&chatID, &usernameA, &usernameB, &createdAt, &sourceLang, &targetLang); err != nil {
			return nil, err
		}
		chats = append(chats, &models.Chat{
			ChatID:     chatID,
			UsernameA:  usernameA,
			UsernameB:  usernameB,
			CreatedAt:  time.Unix(createdAt, 0),
			SourceLang: sourceLang,
			TargetLang: targetLang,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (s *Store) UpdateMessageTranslation(ctx context.Context, messageID string, translatedContent string) error {
	query := `
		UPDATE messages
		SET translated_content = $1
		WHERE message_id = $2
	`
	_, err := s.dbConn.Exec(ctx, query, translatedContent, messageID)
	if err != nil {
		return err
	}

	return nil
}

func scanChatMessage(rs pgx.Row) (*models.ChatMessage, error) {
	var (
		messageID         string
		chatID            string
		senderUsername    string
		receiverUsername  string
		content           string
		translatedContent string
		ts                int64
	)

	if err := rs.Scan(&messageID, &chatID, &senderUsername, &receiverUsername, &content, &translatedContent, &ts); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	return &models.ChatMessage{
		MessageID:         messageID,
		ChatID:            chatID,
		SenderUsername:    senderUsername,
		ReceiverUsername:  receiverUsername,
		Content:           content,
		TranslatedContent: translatedContent,
		Timestamp:         time.Unix(ts, 0),
	}, nil
}
