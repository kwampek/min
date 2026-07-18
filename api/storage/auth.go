package storage

import (
	"api/models"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrLoginExists    = errors.New("login already exists")
	ErrEmailExists    = errors.New("email already exists")
	ErrInvalidSession = errors.New("invalid session")
)

func (s *Storage) CreateUser(user models.User) (int, error) {
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

func (s *Storage) GetUserByLogin(login string) (models.User, error) {
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

	return user, err
}

func (s *Storage) GetUserByID(id int) (models.User, error) {
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

	return user, err
}

// func (s *Storage) IsTokenExists(token string) (bool, error) {
// 	var exists bool
// 	err := s.DB.QueryRow(
// 		`SELECT EXISTS (
//             SELECT 1
//             FROM Messenger.Tokens
//             WHERE token = $1
//               AND expires_at > now()
//         )`, token,
// 	).Scan(&exists)

// 	return exists, err
// }

// func (s *Storage) SaveToken(userID int, token string, expiresAt time.Time) error {
// 	_, err := s.DB.Exec(`
// 		INSERT INTO Messenger.Tokens (user_id, token, expires_at)
// 		VALUES ($1, $2, $3)
// 	`, userID, token, expiresAt)

// 	return err
// }

// func (s *Storage) DeleteToken(userID int, token string) error {
// 	_, err := s.DB.Exec(`
// 		DELETE FROM Messenger.Tokens
// 		WHERE user_id = $1
// 		  AND token = $2
// 	`, userID, token)

// 	return err
// }

// func (s *Storage) ValidateSession(token string) (int, error) {
// 	hash := sha256.Sum256([]byte(token))

// 	var userID int

// 	err := s.DB.QueryRow(`
// 		SELECT user_id
// 		FROM Messenger.Sessions
// 		WHERE token_hash=$1
// 		AND revoked=false
// 		AND expires_at > now()
// 	`,
// 		hex.EncodeToString(hash[:]),
// 	).Scan(&userID)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }

func (s *Storage) ValidateSession(token string) (models.Identifier, error) {
	var userID int
	var sessionID int

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	err := s.DB.QueryRow(`
        UPDATE Messenger.Sessions
        SET last_activity = NOW()
        WHERE token_hash = $1
          AND revoked = FALSE
          AND expires_at > NOW()
        RETURNING 
            session_id,
            user_id,
			token_hash,
            device_name,
            ip_address,
            expires_at
    `, tokenHash).Scan(
		&sessionID,
		&userID,
	)

	if err == sql.ErrNoRows {
		return models.Identifier{}, ErrInvalidSession
	}
	if err != nil {
		return models.Identifier{}, fmt.Errorf("validate session: %w", err)
	}

	return models.Identifier{
		SessionID: sessionID,
		UserID:    userID,
	}, nil
}

func (s *Storage) CreateSession(session models.Session) error {
	_, err := s.DB.Exec(`
		INSERT INTO Messenger.Sessions
		(
			user_id,
			token_hash,
			device_name,
			ip_address,
			expires_at
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`,
		session.UserID,
		session.TokenHash,
		session.DeviceName,
		session.IPAddress,
		session.ExpiresAt,
	)

	return err
}

func (s *Storage) RevokeSession(tokenHash string) error {
	_, err := s.DB.Exec(`
		UPDATE Messenger.Sessions
		SET revoked = TRUE
		WHERE token_hash = $1
	`, tokenHash)

	return err
}
