package repo

import (
	"database/sql"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
)

type UserGroupRepository struct {
	db *sql.DB
}

func NewUserGroupRepository(db *sql.DB) *UserGroupRepository {
	return &UserGroupRepository{db: db}
}

func (g *UserGroupRepository) Create(group dto.UserGroup) (*dto.UserGroup, error) {
	const q = "INSERT INTO main.user_groups(name, is_admin) VALUES($1, $2) RETURNING id"

	var id uint64

	err := g.db.QueryRow(q, group.Name, group.IsAdmin).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &dto.UserGroup{ID: id, Name: group.Name, IsAdmin: group.IsAdmin}, nil
}

func (g *UserGroupRepository) GetByID(id uint64) (*dto.UserGroup, error) {
	const q = "SELECT id, name, is_admin FROM main.user_groups WHERE id = $1"

	var group = dto.UserGroup{}

	row := g.db.QueryRow(q, id)

	err := row.Scan(&group.ID, &group.Name, &group.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &group, nil
}

func (g *UserGroupRepository) Update(group dto.UserGroup) (*dto.UserGroup, error) {
	const q = "UPDATE main.user_groups SET name = $1, is_admin = $2 WHERE id = $3 RETURNING id"

	var id uint64

	err := g.db.QueryRow(q, group.Name, group.IsAdmin, group.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &dto.UserGroup{ID: id, Name: group.Name, IsAdmin: group.IsAdmin}, nil
}

func (g *UserGroupRepository) DeleteByID(id uint64) error {
	const q = "DELETE FROM main.user_groups WHERE id = $1"

	_, err := g.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
