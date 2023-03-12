package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
)

type UserGroupRepository struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewUserGroupRepository(ctx context.Context, db *pgxpool.Pool) *UserGroupRepository {
	return &UserGroupRepository{ctx: ctx, db: db}
}

func (g *UserGroupRepository) Create(group dto.UserGroup) (*dto.UserGroup, error) {
	const q = "INSERT INTO public.user_groups(name, is_admin) VALUES($1, $2) RETURNING id"

	var id uint64

	err := g.db.QueryRow(g.ctx, q, group.Name, group.IsAdmin).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &dto.UserGroup{ID: id, Name: group.Name, IsAdmin: group.IsAdmin}, nil
}

func (g *UserGroupRepository) GetByID(id uint64) (*dto.UserGroup, error) {
	const q = "SELECT id, name, is_admin FROM public.user_groups WHERE id = $1"

	var group = dto.UserGroup{}

	row := g.db.QueryRow(g.ctx, q, id)

	err := row.Scan(&group.ID, &group.Name, &group.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &group, nil
}

func (g *UserGroupRepository) Get(perPage, page int) ([]*dto.UserGroup, error) {
	const q = "SELECT id, name, is_admin FROM public.user_groups  limit $1 offset $2"

	var userGroups = make([]*dto.UserGroup, 0, perPage)

	result, err := g.db.Query(g.ctx, q, perPage, page*perPage)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}

	for result.Next() {
		userGroup := dto.UserGroup{}
		err = result.Scan(&userGroup.ID, &userGroup.Name, &userGroup.IsAdmin)
		if err != nil {
			return nil, err
		}
		userGroups = append(userGroups, &userGroup)
	}

	return userGroups, nil
}

func (g *UserGroupRepository) Update(group dto.UserGroup) (*dto.UserGroup, error) {
	const q = "UPDATE public.user_groups SET name = $1, is_admin = $2 WHERE id = $3 RETURNING id"

	var id uint64

	err := g.db.QueryRow(g.ctx, q, group.Name, group.IsAdmin, group.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &dto.UserGroup{ID: id, Name: group.Name, IsAdmin: group.IsAdmin}, nil
}

func (g *UserGroupRepository) DeleteByID(id uint64) error {
	const q = "DELETE FROM public.user_groups WHERE id = $1"

	_, err := g.db.Exec(g.ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}
