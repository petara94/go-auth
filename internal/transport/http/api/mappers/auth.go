package mappers

import (
	serv_dto "github.com/petara94/go-auth/internal/services/dto"
	"github.com/petara94/go-auth/internal/transport/http/api/dto"
	"time"
)

// LoginReq to Auth
func LoginReqToAuth(loginReq *dto.LoginReq) *serv_dto.Auth {
	if loginReq.TTL == nil {
		return &serv_dto.Auth{
			Login:    loginReq.Login,
			Password: loginReq.Password,
			TTL:      nil,
		}
	}

	var ttl = &time.Time{}

	shift, err := time.ParseDuration(*loginReq.TTL)
	if err != nil {
		ttl = nil
	} else {
		*ttl = time.Now().Add(shift)
	}

	return &serv_dto.Auth{
		Login:    loginReq.Login,
		Password: loginReq.Password,
		TTL:      ttl,
	}
}
