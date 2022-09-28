package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

// ID          uuid.UUID `db:"id" json:"id"`
// Name        string    `db:"name" json:"name"`
// Description string    `db:"description" json:"description"`
// Price       float32   `db:"price" json:"price"`
// CreatedAt   time.Time `db:"created_at" json:"created_at"`
// schema для создания таблиц на RawSQL
var schema = `
-- CREATE TYPE user_role AS ENUM ('customer', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    login TEXT NOT NULL ,
    password TEXT NOT NULL,
    role user_role DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE IF NOT EXISTS items (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price float NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
CREATE TABLE IF NOT EXISTS orders (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    items json NOT NULL,
    price float NOT NULL,
    ordered_by UUID,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_ordered_by FOREIGN KEY(ordered_by) REFERENCES users(id)
);
`

// ConnectDatabase подключается к базе данных
func ConnectDatabase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(0)

	_, err = db.Exec(schema)
	if err != nil {
		log.Println("CreateSchema error:", err)
	}

	return db, err
}
