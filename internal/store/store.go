package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"password-manager/internal/models"
	"password-manager/internal/service_errors"
	"time"
)

const pgDuplicateCode = "23505"

type Store struct {
	store *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{store: db}
}

func (o *Store) CreateUserDB(ctx context.Context, user models.Users) error {
	timestamp := time.Now().Unix()

	stmt, err := o.store.DB.PrepareContext(ctx, createUser) // check the index
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == pgDuplicateCode {
				logrus.Warnf("user with login already exist: %v", err)
				return service_errors.ErrLoginAlreadyExist
			} else {
				logrus.Errorf("unexpected DB error: %v", err)
				return service_errors.ErrWithDB
			}
		} else {
			logrus.Errorf("unhandled error: %v", err)
			return err
		}
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.Login, user.Password, timestamp)
	if err != nil {
		logrus.Errorf("unhandled error: %v", err)
		return err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}

func (o *Store) GetUserPass(ctx context.Context, login string) (string, bool) {

	var hash string

	stmt, err := o.store.DB.PrepareContext(ctx, getUserPass)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
	}

	err = stmt.QueryRowContext(ctx, login).Scan(&hash)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("No rows returned")
			return "", false
		}
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return hash, true
}

func (o *Store) GetUserID(ctx context.Context, login string) (string, error) {
	var UID string

	stmt, err := o.store.DB.PrepareContext(ctx, getUserID)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return "", err
	}

	err = stmt.QueryRowContext(ctx, login).Scan(&UID)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("impossible error: %v", err)
			return "", errors.New("login not found. Impossible")
		}
		logrus.Errorf("unhandled error: %v", err)
		return "", err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return UID, nil
}

func (o *Store) SaveUserPasswordDB(ctx context.Context, req models.Password) error {
	stmt, err := o.store.DB.PrepareContext(ctx, createUserPassword)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, req.UserID, req.Name, req.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == pgDuplicateCode {
				logrus.Info("name already exist on this user: %v", err)
				return service_errors.ErrNameAlreadyExist
			}
		}
		logrus.Errorf("unhandled error: %v", err)
		return err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return nil
}

func (o *Store) GetUserPasswordDB(ctx context.Context, name, UID string) (models.Password, error) {

	var res models.Password
	stmt, err := o.store.DB.PrepareContext(ctx, getUserSavedPassword)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return models.Password{}, err
	}

	err = stmt.QueryRowContext(ctx, name, UID).Scan(&res.Name, &res.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("Passwords with name doesn't exist: %v", err)
			return models.Password{}, service_errors.ErrPasswordNotFound
		}
		logrus.Errorf("unhandled error: %v", err)
		return models.Password{}, err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return res, nil
}

func (o *Store) UpdateUserSavedPasswordDB(ctx context.Context, req models.NewPassword) error {
	stmt, err := o.store.DB.PrepareContext(ctx, updateUserPassword)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, req.NewName, req.NewPassword, req.OldName, req.UserID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == pgDuplicateCode {
				logrus.Info("name already exist on this user: %v", err)
				return service_errors.ErrNameAlreadyExist
			}
		}
		logrus.Errorf("unhandled error: %v", err)
		return err
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}
	return nil
}

func (o *Store) DeleteUserSavedPasswordDB(ctx context.Context, name, UID string) error {

	stmt, err := o.store.DB.PrepareContext(ctx, deleteUserPassword)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return err
	}

	_, err = stmt.ExecContext(ctx, name, UID)
	if err != nil {
		logrus.Errorf("can not delete the user's password")
		return service_errors.ErrWithDB
	}

	if err := stmt.Close(); err != nil {
		logrus.Warnf("attention error closing statment: %v", err)
	}

	return nil
}

func (o *Store) GetAllUserPasswordDB(ctx context.Context, UID string) ([]models.PasswordName, error) {

	var passwords []models.PasswordName
	var pass models.PasswordName

	stmt, err := o.store.DB.PrepareContext(ctx, getAllUserSavedPassword)

	if err != nil {
		logrus.Errorf("error with stmt: %v", err)
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, UID)

	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("Passwords for current user not found: %v", err)
			return nil, service_errors.ErrPasswordNotFound
		}
		logrus.Errorf("unhandled error: %v", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pass.Name)
		if err != nil {
			return nil, err
		}
		passwords = append(passwords, pass)
	}

	return passwords, nil
}
