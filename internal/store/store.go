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
			if pqErr.Code == "23505" {
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

	err = stmt.QueryRowContext(ctx, name, UID).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Info("Passwords with name doesn't exist: %v", err)
			return models.Password{}, service_errors.ErrPasswordNotFound
		}
		logrus.Errorf("unhandled error: %v", err)
		return models.Password{}, err
	}

	return res, nil
}
