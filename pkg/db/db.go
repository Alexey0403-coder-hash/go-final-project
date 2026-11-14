package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(30) NOT NULL DEFAULT '',
	comment TEXT DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX date_index
ON scheduler (date);
`

var db *sql.DB

// Init запускает пул соединения с БД и, если это необходимо, создает БД
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	var install bool
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		return err
	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install == true {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
