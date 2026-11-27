package store

import (
	"database/sql"
	"goLogMyTime/common/globals"
	"goLogMyTime/utils"
	"sync"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var once sync.Once

const FILE_PATH = globals.DATA_DIR + "db.db"

type DbStore struct {
	dbInstance *sql.DB
}

func Init(s *DbStore) {
	once.Do(func() {
		db, err := sql.Open("sqlite3", FILE_PATH)
		utils.Assert(err == nil, "failed to start sqlite")
		defer db.Close()

		createTaskDefsTable := `CREATE TABLE IF NOT EXISTS taskDefinitions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ticket TEXT NOT NULL,
			info TEXT NOT NULL,
			isFavorite BOOLEAN NOT NULL DEFAULT FALSE,
			timestampOfLastActivity DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
		_, err = db.Exec(createTaskDefsTable)
		utils.Assert(err == nil, "failed to create taskDefs table")

		createTaskRecordsTable := `CREATE TABLE IF NOT EXISTS taskRecords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ticket TEXT NOT NULL,
			timestampOfStart DATETIME DEFAULT CURRENT_TIMESTAMP
			isActive BOOLEAN NOT NULL DEFAULT TRUE,
			timestampOfEnd DATETIME DEFAULT 0
		);`
		_, err = db.Exec(createTaskRecordsTable)
		utils.Assert(err == nil, "failed to create taskRecords table")

		s.dbInstance = db
	})
}
