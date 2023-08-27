package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"password-manager/internal/config"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
	"time"
)

type StoreInterface interface {
	CreateUserDB(ctx context.Context, user models.Users) error
	GetUserPass(ctx context.Context, login string) (string, bool)
	GetUserID(ctx context.Context, login string) (string, error)
	SaveUserPasswordDB(ctx context.Context, req models.Password) error
	GetUserPasswordDB(ctx context.Context, name, UID string) (models.Password, error)
	UpdateUserSavedPasswordDB(ctx context.Context, req models.NewPassword) error
	DeleteUserSavedPasswordDB(ctx context.Context, name, UID string) error
	GetAllUserPasswordDB(ctx context.Context, UID string) ([]models.PasswordName, error)
}

type Service struct {
	store  StoreInterface
	config config.Config
}

func NewService(s StoreInterface, config *config.Config) *Service {
	return &Service{
		store:  s,
		config: *config,
	}
}

func (s *Service) CreateUser(ctx context.Context, req models.Users) (string, error) {
	if req.Login == "" || req.Password == "" {
		return "", service_errors.ErrBadRequest
	}
	req.Password = generatePasswordHash(req.Password, s.config.PasswordSecretValue)
	req.ID = generateUID()

	if err := s.store.CreateUserDB(ctx, req); err != nil { // add index to check the login
		return "", service_errors.ErrLoginAlreadyExist
	}
	token := NewToken(req.ID)

	return token.SignedString([]byte(s.config.SecretValue))
}

func (s *Service) GenerateUserToken(ctx context.Context, req models.Users) (string, error) {
	passwordHash, exist := s.store.GetUserPass(ctx, req.Login)
	if !exist {
		return "", service_errors.ErrInvalidLoginOrPass
	}

	isPassValid := comparePasswordHash(passwordHash, req.Password, s.config.PasswordSecretValue)
	if isPassValid {
		UID, err := s.store.GetUserID(ctx, req.Login)
		if err != nil {
			return "", errors.New("error while getting UID: %v")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenTTL).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UID,
		})

		return token.SignedString([]byte(s.config.SecretValue))
	}
	return "", service_errors.ErrInvalidLoginOrPass
}

func (s *Service) SaveUserPassword(ctx context.Context, req models.Password) error {
	if req.Name == "" {

		return service_errors.ErrEmptyNameOrPassword
	}
	err := s.store.SaveUserPasswordDB(ctx, req)
	if err != nil {
		return service_errors.ErrWithDB
	}
	return nil
}

func (s *Service) GetUserPassword(ctx context.Context, name, UID string) (models.Password, error) {
	if name == "" {
		return models.Password{}, service_errors.ErrBadRequest
	}

	res, err := s.store.GetUserPasswordDB(ctx, name, UID)

	if err != nil {
		return models.Password{}, err
	}
	return res, nil
}

func (s *Service) UpdateUserSavedPassword(ctx context.Context, req models.NewPassword) error {
	if req.NewName == "" || req.NewPassword == "" || req.OldName == "" {
		return service_errors.ErrEmptyNameOrPassword
	} // подумать над проверкой уже имеющихся данных

	err := s.store.UpdateUserSavedPasswordDB(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteUserSavedPassword(ctx context.Context, name, UID string) error {
	if name == "" {
		return service_errors.ErrBadRequest
	}

	err := s.store.DeleteUserSavedPasswordDB(ctx, name, UID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetAllUserPasswords(ctx context.Context, UID string) ([]models.PasswordName, error) {

	pass, err := s.store.GetAllUserPasswordDB(ctx, UID)
	if err != nil {
		return nil, err
	}
	return pass, nil
}
