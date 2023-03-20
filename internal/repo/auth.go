package repo

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petara94/go-auth/internal/services/dto"
)

type SessionStore struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewSessionStore(ctx context.Context, db *pgxpool.Pool) *SessionStore {
	return &SessionStore{ctx: ctx, db: db}
}

func (s *SessionStore) Create(session dto.Session) (*dto.Session, error) {
	const q = "INSERT INTO public.sessions (token, user_id, expr) VALUES ($1, $2, $3) RETURNING token, user_id, expr"

	result := s.db.QueryRow(s.ctx, q, session.Token, session.UserID, session.Expr)

	var createdSession dto.Session
	err := result.Scan(&createdSession.Token, &createdSession.UserID, &createdSession.Expr)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &createdSession, nil
}

func (s *SessionStore) GetByToken(token string) (*dto.Session, error) {
	const q = "SELECT token, user_id, expr FROM public.sessions WHERE token = $1"

	var session dto.Session

	result := s.db.QueryRow(s.ctx, q, token)

	err := result.Scan(&session.Token, &session.UserID, &session.Expr)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (s *SessionStore) DeleteByToken(token string) error {
	const q = "DELETE FROM public.sessions WHERE token = $1"

	_, err := s.db.Exec(s.ctx, q, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionStore) DeleteByUserID(userID uint64) error {
	const q = "DELETE FROM public.sessions WHERE user_id = $1"

	_, err := s.db.Exec(s.ctx, q, userID)
	if err != nil {
		return err
	}

	return nil
}
