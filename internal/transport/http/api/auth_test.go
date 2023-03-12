package api

import (
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func getAuthServiseMock(t *testing.T) AuthService {
	a := mocks.NewAuthService(t)
	e := a.EXPECT()

	sess1 := &dto.Session{
		Token:  "1",
		UserID: 1,
		Expr:   nil,
	}

	e.Get("1").Return(sess1, nil)

	e.Login(dto.Auth{
		Login:    "1",
		Password: "1",
	}).Return(sess1, nil)

	e.Logout(*sess1).Return(nil)

	return a
}

func TestLoginLogout(t *testing.T) {
	authServise := getAuthServiseMock(t)

	a := dto.Auth{
		Login:    "1",
		Password: "1",
	}

	session1, err := authServise.Login(a)
	require.NoError(t, err)

	got, err := authServise.Get(session1.Token)
	require.NoError(t, err)

	require.Equal(t, *got, *session1)

	err = authServise.Logout(*session1)
	require.NoError(t, err)
}
