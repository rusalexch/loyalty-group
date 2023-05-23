package auth

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Mux         *chi.Mux
	UserService userService
	jwtSecret   string
}

type authModule struct {
	mux         *chi.Mux
	userService userService
	jwtSecret   string
}

func New(conf Config) *authModule {
	module := &authModule{
		mux:         conf.Mux,
		userService: conf.UserService,
		jwtSecret:   conf.jwtSecret,
	}
	module.init()

	return module
}

func (am *authModule) CheckToken(ctx context.Context, authToken string) (app.User, error) {
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return am.jwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user"]; ok {
			user, err := am.userService.FindByID(ctx, userID.(int))
			if err != nil {
				return app.User{}, err
			}

			return user, nil
		}
	}

	return app.User{}, err
}

func (am *authModule) init() {
	am.mux.Post(registerPattern, am.register)
	am.mux.Post(loginPattern, am.login)
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
	t := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user": userID,
	})
	return t.SignedString(am.jwtSecret)
}

func (am *authModule) signup(ctx context.Context, signup signup) (string, error) {
	user, err := am.userService.FundByLogin(ctx, signup.Login)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		return "", err
	}
	if user.Login == signup.Login {
		return "", errLoginAlreadyExist
	}

	hash, err := hash(signup.Password)
	if err != nil {
		return "", err
	}
	user, err = am.userService.Create(ctx, app.CreateUser{
		Login:    signup.Login,
		Password: hash,
	})
	if err != nil {
		return "", err
	}

	return am.token(user.ID)
}

func (am *authModule) signin(ctx context.Context, login login) (string, error) {
	user, err := am.userService.FundByLogin(ctx, login.Login)
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
