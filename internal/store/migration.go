package store

import "fmt"

func (o *Store) MakeMigration() error {
	if _, err := o.store.DB.Exec(createUserTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createPasswordTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createUserKeyTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createFileTable); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createPasswordForeignKey); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}
	if _, err := o.store.DB.Exec(createFileForeignKey); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createLoginIndex); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createIDNameIndex); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	if _, err := o.store.DB.Exec(createIDNameFileIndex); err != nil {
		return fmt.Errorf("error while run migrations %v", err)
	}

	return nil
}
