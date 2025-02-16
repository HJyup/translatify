package store

import (
	"context"
	"errors"
	models "github.com/HJyup/translatify-user/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	"github.com/HJyup/translatify-common/utils"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) CreateUser(ctx context.Context, username, email, password string) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`
	var userId string
	err := s.dbConn.QueryRow(ctx, query, username, email, password).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, status.Error(codes.AlreadyExists, "user with this email or username is already registered")
		}
		return nil, err
	}

	return &models.User{
		UserId:    userId,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}, nil
}

func (s *Store) GetUser(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT user_id, username, email, password, created_at
		FROM users
		WHERE username = $1
	`
	row := s.dbConn.QueryRow(ctx, query, username)
	return scanUser(row)
}

func (s *Store) DeleteUser(ctx context.Context, userId string) (bool, error) {
	query := `
		DELETE FROM users
		WHERE user_id = $1
	`
	tag, err := s.dbConn.Exec(ctx, query, userId)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() == 0 {
		return false, errors.New("user not found")
	}
	return true, nil
}

func (s *Store) ListUsers(ctx context.Context, limit int, paginationToken string) ([]*models.User, string, error) {
	var effectiveSince int64 = 0
	if paginationToken != "" {
		if tokenTs, err := strconv.ParseInt(paginationToken, 10, 64); err == nil && tokenTs > effectiveSince {
			effectiveSince = tokenTs
		}
	}

	query := `
		SELECT user_id, username, email, password, created_at
		FROM users
		WHERE created_at > to_timestamp($1)
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := s.dbConn.Query(ctx, query, effectiveSince, limit+1)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, "", err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(users) > limit {
		nextPageToken = strconv.FormatInt(users[limit].CreatedAt.Unix(), 10)
		users = users[:limit]
	}

	return users, nextPageToken, nil
}

func scanUser(rs utils.RowScanner) (*models.User, error) {
	var (
		userId    string
		username  string
		email     string
		password  string
		createdAt time.Time
	)
	if err := rs.Scan(&userId, &username, &email, &password, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &models.User{
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}
