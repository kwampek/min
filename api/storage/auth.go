package storage

import (
	"api/models"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrLoginExists = errors.New("")
	ErrEmailExists = errors.New("")
)

func (s *Storage) CreateUser(user *models.User) (int, error) {
	var id int

	err := s.DB.QueryRow(`
		INSERT INTO Messenger.Users (
			login,
			password_hash,
			email,
			phone_number,
			avatar_link
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`,
		user.Login,
		user.PasswordHash,
		user.Email,
		user.PhoneNumber,
		user.AvatarLink,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			switch pgErr.ConstraintName {
			case "idx_users_login_unique":
				return -1, ErrLoginExists
			case "idx_users_email_unique":
				return -1, ErrEmailExists
			}
		}
	}

	return id, nil
}

func (s *Storage) GetUserByLogin(login string) (*models.User, error) {
	var user models.User

	err := s.DB.QueryRow(`
		SELECT
			user_id,
			login,
			password_hash,
			email,
			phone_number,
			avatar_link,
			created_at,
			search_privacy
		FROM Messenger.Users
		WHERE LOWER(login) = LOWER($1)
	`, login).Scan(
		&user.UserID,
		&user.Login,
		&user.PasswordHash,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarLink,
		&user.CreatedAt,
		&user.SearchPrivacy,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetUserByID(id int) (*models.User, error) {
	var user models.User

	err := s.DB.QueryRow(`
		SELECT
			user_id,
			login,
			password_hash,
			email,
			phone_number,
			avatar_link,
			created_at,
			search_privacy
		FROM Messenger.Users
		WHERE user_id = $1
	`, id).Scan(
		&user.UserID,
		&user.Login,
		&user.PasswordHash,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarLink,
		&user.CreatedAt,
		&user.SearchPrivacy,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) SaveToken(userID int, token string, expiresAt time.Time) error {
	_, err := s.DB.Exec(`
		INSERT INTO Messenger.Tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`, userID, token, expiresAt)

	return err
}

func (s *Storage) DeleteToken(userID int, token string) error {
	_, err := s.DB.Exec(`
		DELETE FROM Messenger.Tokens
		WHERE user_id = $1
		  AND token = $2
	`, userID, token)

	return err
}
