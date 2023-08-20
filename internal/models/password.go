package models

type Password struct {
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"data"`
	UserID   string `json:"-"`
}

type NewPassword struct {
	NewName     string `json:"new_name" db:"name"`
	NewPassword string `json:"new_password" db:"data"`
	OldName     string `json:"old_name"`
	UserID      string `json:"-"`
}
