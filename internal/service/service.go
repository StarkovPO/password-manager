package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	cipher_client "password-manager/internal/cipher"
	"password-manager/internal/config"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
	"time"
)

// StoreInterface is interface for store module
type StoreInterface interface {
	CreateUserDB(ctx context.Context, user models.Users) error
	GetUserPass(ctx context.Context, login string) (string, bool)
	GetUserID(ctx context.Context, login string) (string, error)
	SaveUserPasswordDB(ctx context.Context, req models.Password) error
	GetUserPasswordDB(ctx context.Context, name, UID string) (models.Password, error)
	UpdateUserSavedPasswordDB(ctx context.Context, req models.NewPassword) error
	DeleteUserSavedPasswordDB(ctx context.Context, name, UID string) error
	GetAllUserPasswordDB(ctx context.Context, UID string) ([]models.PasswordName, error)
	SaveUserKey(ctx context.Context, UID string, key string) error
	GetUserKey(ctx context.Context, UID string) (string, error)
}

// Service struct for service
type Service struct {
	store  StoreInterface
	config config.Config
}

// NewService create new service
func NewService(s StoreInterface, config *config.Config) *Service {
	return &Service{
		store:  s,
		config: *config,
	}
}

/*
CreateUser create new user. Return token and error.
Accept the object type of models.Users
Generate password hash for user and save in DB
Generate UID for user and save in DB
*/
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

/*
	GenerateUserToken generate user token. Used for auth users

Hash users password and compare with DB hash password
Gets user ID and generate token
Accept the object type of models.Users
*/
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

/*
SaveUserPassword save user password from request
Accept the object type of models.Password
*/
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

/*
GetUserPassword get user password from DB
Accept the object type of models.Password
Return encrypted user password
*/
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

/*
UpdateUserSavedPassword update user password from DB
Accept the object type of models.NewPassword
Return error
*/
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

/*
DeleteUserSavedPassword delete user password from DB
Accept the object type of models.NewPassword
Return error
*/
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

/*
GetAllUserPasswords get all user password names from DB
Accept the object type of models.Password
Return array of users password names
*/
func (s *Service) GetAllUserPasswords(ctx context.Context, UID string) ([]models.PasswordName, error) {

	pass, err := s.store.GetAllUserPasswordDB(ctx, UID)
	if err != nil {
		return nil, err
	}
	return pass, nil
}

func (s *Service) GetUserKey(ctx context.Context, UID string) (string, error) {
	key, err := s.store.GetUserKey(ctx, UID)
	if err != nil {
		return "", err
	}
	if key != "" {
		decryptedKey, err := cipher_client.Decrypt(key, []byte(s.config.PasswordSecretValue))
		if err != nil {
			return "", err
		}

		return decryptedKey, nil
	}
	return "", nil
}

func (s *Service) SaveUserKey(ctx context.Context, UID string, key string) error {

	encryptedKey, err := cipher_client.Encrypt(key, []byte(s.config.PasswordSecretValue))

	if err != nil {
		return err
	}

	err = s.store.SaveUserKey(ctx, UID, encryptedKey)

	if err != nil {
		return err
	}
	return nil
}
