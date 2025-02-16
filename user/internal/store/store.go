package store

import (
	"context"
	"errors"
	"github.com/HJyup/translatify-common/utils"
	"strconv"
	"time"

	pb "github.com/HJyup/translatify-common/api"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

type Store struct {
	dbConn *pgx.Conn
}

func NewStore(dbConn *pgx.Conn) *Store {
	return &Store{dbConn: dbConn}
}

func (s *Store) CreateUser(ctx context.Context, username, email, fullName, password string) (*pb.User, error) {
	query := `
		INSERT INTO users
			(username, email, full_name, password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`
	ts := time.Now().Unix()
	var userId string
	err := s.dbConn.QueryRow(ctx, query, username, email, fullName, password, ts).Scan(&userId)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		UserId:    userId,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		Password:  password,
		CreatedAt: ts,
	}, nil
}

// GetUser retrieves a user by user_id.
func (s *Store) GetUser(ctx context.Context, userId string) (*pb.User, error) {
	query := `
		SELECT user_id, username, email, full_name, password, created_at
		FROM users
		WHERE user_id = $1
	`
	row := s.dbConn.QueryRow(ctx, query, userId)
	return scanUser(row)
}

func (s *Store) UpdateUser(ctx context.Context, username, email, fullName, password string) (*pb.User, error) {
	query := `
		UPDATE users
		SET email = $2,
		    full_name = $3,
		    password = $4,
		    updated_at = $5
		WHERE username = $1
		RETURNING user_id, created_at
	`
	var userId string
	var createdAt int64
	err := s.dbConn.QueryRow(ctx, query, username, email, fullName, password).Scan(&userId, &createdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &pb.User{
		UserId:    userId,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		Password:  password,
		CreatedAt: createdAt,
	}, nil
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
		return false, ErrUserNotFound
	}
	return true, nil
}

func (s *Store) ListUsers(ctx context.Context, limit int, paginationToken string) ([]*pb.User, string, error) {
	var effectiveSince int64 = 0
	if paginationToken != "" {
		if tokenTs, err := strconv.ParseInt(paginationToken, 10, 64); err == nil && tokenTs > effectiveSince {
			effectiveSince = tokenTs
		}
	}

	query := `
		SELECT user_id, username, email, full_name, password, created_at
		FROM users
		WHERE created_at > $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := s.dbConn.Query(ctx, query, effectiveSince, limit+1)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	users := make([]*pb.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, "", err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}

	var nextPageToken string
	if len(users) > limit {
		nextPageToken = strconv.FormatInt(users[limit].CreatedAt, 10)
		users = users[:limit]
	}

	return users, nextPageToken, nil
}

func scanUser(rs utils.RowScanner) (*pb.User, error) {
	var (
		userId    string
		username  string
		email     string
		fullName  string
		password  string
		createdAt int64
	)
	if err := rs.Scan(&userId, &username, &email, &fullName, &password, &createdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &pb.User{
		UserId:    userId,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		Password:  password,
		CreatedAt: createdAt,
	}, nil
}
