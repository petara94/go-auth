package api

import (
	"errors"
	"github.com/google/uuid"
	"github.com/petara94/go-auth/internal/services/pkg"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func getUserServiceMock() UserService {
	u := mocks.UserService{}
	e := u.EXPECT()

	var (
		userId uint64 = 1
	)

	e.Get(-1, 1).Return(nil, errors.New("mock"))
	e.Get(1, -1).Return(nil, errors.New("mock"))
	e.Get(1, 1).Return([]*dto.User{
		{
			ID:          1,
			UserGroupID: nil,
			Login:       uuid.NewString(),
			Password:    uuid.NewString(),
		},
		{
			ID:          2,
			UserGroupID: &userId,
			Login:       uuid.NewString(),
			Password:    uuid.NewString(),
		},
	}, nil)

	e.Create(dto.User{
		Login:    "111",
		Password: "111",
	}).Return(&dto.User{
		ID:          1,
		Login:       "111",
		Password:    pkg.HashPassword("111"),
		UserGroupID: nil,
	}, nil)

	return &u
}

func TestUserCreate(t *testing.T) {
	userService := getUserServiceMock()

	user, err := userService.Create(dto.User{Login: "111", Password: "111"})
	require.NoError(t, err)

	require.Equal(t, *user, dto.User{ID: 1, Login: "111", Password: pkg.HashPassword("111")})
}
