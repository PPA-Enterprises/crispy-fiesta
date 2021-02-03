package users

import (
	"context"
	passwordUtils "internal/common/password"
	"internal/common/errors"
)

func authenticate(ctx context.Context, credentials loginUserCommand) (string, *errors.ResponseError){
	user, err := UserByEmail(ctx, credentials.Email); if err != nil {
		//email doesnt exist
		return string(""), errors.EmailDoesNotExistError()
	}

	if user.IsDeleted {
		// User was deleted
		return string(""), errors.EmailDoesNotExistError()
	}

	ok, err := passwordUtils.VerifyPassword(credentials.Password, user.Password)
	if !ok {
		//passwords dont match
		return string(""), errors.InvalidCredentials()
	}
	//get jwt
	jwt, err := user.jwt(); if err != nil {
		//failed to create jwt
		return string(""), errors.JwtError(err)
	}
	return jwt, nil
}

func EmailExists(ctx context.Context, email string) bool {
	_, err := UserByEmail(ctx, email); if err != nil {
		return false
	}
	return true
}
