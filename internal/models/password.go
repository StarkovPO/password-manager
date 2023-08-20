package models

type Password struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	UserID   string `json:"-"`
}
