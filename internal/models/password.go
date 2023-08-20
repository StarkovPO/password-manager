package models

type Password struct {
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"data"`
	UserID   string `json:"-"`
}
