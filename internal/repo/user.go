package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
)

type UserStore struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewUserStore(ctx context.Context, db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db, ctx: ctx}
}

func (u *UserStore) Create(user dto.User) (uint64, error) {
	const q = "INSERT INTO public.users(login, password) VALUES($1, $2) RETURNING id"

	var id uint64

	err := u.db.QueryRow(
		u.ctx, q,
		user.Login,
		[]byte(user.Password),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("sql query error: %w", err)
	}

	return id, nil
}

func (u *UserStore) LinkToGroup(id, groupId uint64) error {
	const q = "UPDATE public.users set user_group_id = $1 where id = $2"

	_, err := u.db.Query(u.ctx, q, id, groupId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) GetByID(id uint64) (*dto.User, error) {
	const q = "SELECT id, login, password, user_group_id from public.users where id = $1"

	var (
		user = dto.User{}
		pass []byte
	)

	result, err := u.db.Query(u.ctx, q, id)
	if err != nil {
		return nil, err
	}
	if !result.Next() {
		return nil, ErrNotFound
	}

	err = result.Scan(&user.ID, &user.Login, &pass, &user.UserGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	user.Password = string(pass)

	return &user, nil
}

func (u *UserStore) Get(perPage, page int) ([]*dto.User, error) {
	const q = "SELECT id, login, password, user_group_id from public.users limit $1 offset $2"

	var users = make([]*dto.User, 0, perPage)

	result, err := u.db.Query(u.ctx, q, perPage, page*perPage)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}

	for result.Next() {
		user := dto.User{}
		err = result.Scan(&user.ID, &user.Login, &user.Password, &user.UserGroupID)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *UserStore) GetByLogin(login string) (*dto.User, error) {
	const q = "SELECT id, login, password, user_group_id FROM public.users WHERE login = $1::text"

	var user = dto.User{}

	row := u.db.QueryRow(u.ctx, q, login)

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.UserGroupID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) Update(user dto.User) (*dto.User, error) {
	const q = "UPDATE public.users SET login = $1, password = $2, user_group_id = $3 WHERE id = $4 RETURNING id, login, password, user_group_id"

	var updatedUser = dto.User{}

	row := u.db.QueryRow(u.ctx, q, user.Login, user.Password, user.UserGroupID, user.ID)

	err := row.Scan(&updatedUser.ID, &updatedUser.Login, &updatedUser.Password, &updatedUser.UserGroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
