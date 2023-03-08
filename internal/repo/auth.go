package repo

import (
	"database/sql"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) Create(session dto.Session) (*dto.Session, error) {
	const q = "INSERT INTO main.sessions (token, user_id, expr) VALUES ($1, $2, $3) RETURNING token, user_id, expr"

	result := s.db.QueryRow(q, session.Token, session.UserID, session.Expr)

	var createdSession dto.Session
	err := result.Scan(&createdSession.Token, &createdSession.UserID, &createdSession.Expr)
	if err != nil {
		return nil, err
	}

	return &createdSession, nil
}

func (s *SessionStore) GetByToken(token string) (*dto.Session, error) {
	const q = "SELECT token, user_id, expr FROM main.sessions WHERE token = $1"

	var session dto.Session

	result := s.db.QueryRow(q, token)

	err := result.Scan(&session.Token, &session.UserID, &session.Expr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) DeleteByToken(token string) error {
	const q = "DELETE FROM main.sessions WHERE token = $1"

	_, err := s.db.Exec(q, token)
	if err != nil {
		return err
	}

	return nil
}
