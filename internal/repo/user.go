package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petara94/go-auth/internal/services/dto"
)

type UserStore struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewUserStore(ctx context.Context, db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db, ctx: ctx}
}

func (u *UserStore) Create(user dto.User) (uint64, error) {
	const q = "INSERT INTO public.users(login, password, is_admin, is_blocked, check_password) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id uint64

	err := u.db.QueryRow(u.ctx, q, user.Login, user.Password, user.IsAdmin, user.IsBlocked, user.CheckPassword).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("sql query error: %w", err)
	}

	return id, nil
}

func (u *UserStore) GetByID(id uint64) (*dto.User, error) {
	const q = "SELECT id, login, password, is_admin, is_blocked, check_password FROM public.users WHERE id = $1"

	var (
		user = dto.User{}
	)

	result, err := u.db.Query(u.ctx, q, id)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, ErrNotFound
	}

	err = result.Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.IsBlocked, &user.CheckPassword,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) GetWithPagination(perPage, page int) ([]*dto.User, error) {
	const q = "SELECT id, login, password, is_admin, is_blocked, check_password FROM public.users LIMIT $1 OFFSET $2"

	var (
		users = make([]*dto.User, 0, perPage)
		user  = dto.User{}
	)

	result, err := u.db.Query(u.ctx, q, perPage, perPage*(page-1))
	if err != nil {
		return nil, err
	}

	for result.Next() {
		err = result.Scan(
			&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.IsBlocked, &user.CheckPassword,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, ErrNotFound
			}
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *UserStore) GetByLogin(login string) (*dto.User, error) {
	const q = "SELECT id, login, password, is_admin, is_blocked, check_password FROM public.users WHERE login = $1"

	var (
		user = dto.User{}
	)

	result, err := u.db.Query(u.ctx, q, login)
	if err != nil {
		return nil, err
	}

	if !result.Next() {
		return nil, ErrNotFound
	}

	err = result.Scan(
		&user.ID, &user.Login, &user.Password, &user.IsAdmin, &user.IsBlocked, &user.CheckPassword,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) Update(user dto.User) (*dto.User, error) {
	const q = "UPDATE public.users SET login = $1, password = $2, is_admin = $3, is_blocked = $4, check_password = $5 WHERE id = $6 RETURNING id, login, password, is_admin, is_blocked, check_password"

	var (
		updatedUser = dto.User{}
	)

	err := u.db.QueryRow(u.ctx, q, user.Login, user.Password, user.IsAdmin, user.IsBlocked, user.CheckPassword, user.ID).Scan(
		&updatedUser.ID, &updatedUser.Login, &updatedUser.Password, &updatedUser.IsAdmin, &updatedUser.IsBlocked, &updatedUser.CheckPassword,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (u *UserStore) DeleteByID(id uint64) error {
	const q = "DELETE FROM public.users WHERE id = $1"

	_, err := u.db.Exec(u.ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}
