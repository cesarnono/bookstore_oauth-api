package rest

import (
	"encoding/json"
	"golang-restclient/rest"
	"time"

	"github.com/cesarnono/bookstore_oauth-api/src/domain/users"
	"github.com/cesarnono/bookstore_oauth-api/src/util/rest_errors"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8082",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErrr rest_errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErrr) //rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErrr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login response")
	}
	return &user, nil
}
