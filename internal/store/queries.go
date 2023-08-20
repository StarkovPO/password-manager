package store

const (
	createUserTable = `CREATE TABLE IF NOT EXISTS "users" (
        "primary_id" SERIAL PRIMARY KEY,
        "id" varchar(36) UNIQUE,
        "login" varchar(255) UNIQUE,
        "password_hash" varchar(255),
        "created_at" timestamp NOT NULL
    )`

	createPasswordTable = `CREATE TABLE IF NOT EXISTS "passwords" (
	"primary_id" SERIAL PRIMARY KEY,
	"user_id" varchar(36),
	"name" varchar(255) NOT NULL ,
	"data" varchar(1024) NOT NULL 
)`
	createFileTable = `
		CREATE TABLE IF NOT EXISTS "files" (
			"primary_id" SERIAL PRIMARY KEY,
			"user_id" varchar(36),
			"name" varchar(255) NOT NULL ,
			"data" varchar(2048) NOT NULL 
)
`

	createPasswordForeignKey = `ALTER TABLE "passwords" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


 `
	createFileForeignKey = `ALTER TABLE "files" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id")`

	createLoginIndex = `CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex ON public.users (login)`

	createIDNameIndex = `CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_id_name ON public.passwords (user_id, name)`

	createIDNameFileIndex = `CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_id_name_files ON public.files (user_id, name)`

	createUser = `
        INSERT INTO users (id, login, password_hash, created_at)
        VALUES ($1, $2, $3, to_timestamp($4))
    `

	getUserPass = `SELECT password_hash FROM users WHERE login = $1 LIMIT 1`

	getUserID = `SELECT id FROM users WHERE login = $1`

	createUserPassword = `INSERT INTO "passwords" (user_id, name, data) VALUES ($1, $2, $3)`

	getUserSavedPassword = `SELECT name, data FROM passwords WHERE name = $1 AND user_id = $2 limit 1`

	updateUserPassword = `UPDATE passwords SET name = $1, data = $2 WHERE name = $3 and user_id = $4`
)
