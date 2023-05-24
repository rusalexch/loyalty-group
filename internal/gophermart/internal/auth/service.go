package auth

import (
	"context"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
	"golang.org/x/crypto/bcrypt"
)

func (am *authModule) checkToken(ctx context.Context, authToken string) (app.User, error) {
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(am.jwtSecret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user"]; ok {
			user, err := am.userService.FindByID(ctx, int(userID.(float64)))
			if err != nil {
				return app.User{}, err
			}

			return user, nil
		}
	}

	return app.User{}, err
}

func hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash), err
}

func validate(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (am *authModule) token(userID int) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": userID,
	})
	return t.SignedString([]byte(am.jwtSecret))
}

func (am *authModule) signup(ctx context.Context, signup signup) (string, error) {
	user, err := am.userService.FindByLogin(ctx, signup.Login)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		log.Println("auth > signup > can't find user")
		return "", err
	}
	if user.Login == signup.Login {
		return "", errLoginAlreadyExist
	}

	hash, err := hash(signup.Password)
	if err != nil {
		log.Println("auth > signup > can't generate hash")
		return "", err
	}
	user, err = am.userService.Create(ctx, app.CreateUser{
		Login:    signup.Login,
		Password: hash,
	})
	if err != nil {
		log.Println("auth > signup > can't create user")
		return "", err
	}

	return am.token(user.ID)
}

func (am *authModule) signin(ctx context.Context, login login) (string, error) {
	user, err := am.userService.FindByLogin(ctx, login.Login)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			return "", errUnauthorized
		}
		return "", err
	}
	isValid, _ := validate(login.Password, user.Password)
	if !isValid {
		return "", errUnauthorized
	}

	return am.token(user.ID)
}
