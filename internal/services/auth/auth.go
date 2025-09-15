package auth

import (
	"context"
	"fmt"
	"github.com/evgeniySeleznev/gRPC-auth/internal/domain/models"
	"github.com/evgeniySeleznev/gRPC-auth/internal/lib/logger/sl"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

//сервисный слой — для бизнес-логики и взаимодействия с БД (хендлер напрямую не сходит)
//модели для взаимодействия слоёв — в отдельном пакете? чтобы при импорте модели в хендлере —
//он не получал доступ к базе (если бы модель юзера лежала в ней)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Get checks if user with given credentials exists in the system
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error
func (a *Auth) Get(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Get"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		//1:53:30
	}
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error
func (a *Auth) RegisterNewUser(ctx context.Context, email, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	//создаёт пароль, генерит 10 раз соль, возвращает хэш. потом сравниваем через Compare
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
