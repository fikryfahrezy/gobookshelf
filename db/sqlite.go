package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var sqliteDB *sql.DB

func initSqlite(fd string) (*sql.DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fn := filepath.Join(wd, fd)

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitSqliteDB(fd string) (*sql.DB, error) {
	db, err := initSqlite(fd)
	if err != nil {
		return nil, err
	}

	sqliteDB = db

	return db, nil
}

func InitSqliteTestDB(fd string) (*sql.DB, error) {
	db, err := initSqlite(fd)
	if err != nil {
		return nil, err
	}

	sqliteDB = db

	return db, nil
}

func RemoveSqliteTestDB(fd string) error {
	sqliteDB.Close()
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(wd)

	err = os.Remove(fd)
	return err
}

func GetSqliteDB() *sql.DB {
	return sqliteDB
}

func MigrateSqliteDB() {
	tables := []string{
		`
			CREATE TABLE IF NOT EXISTS user_forgot_pass (
				id TEXT,
				email TEXT,
				code TEXT,
				is_claimed INTEGER 
			);
		`,
	}

	db := GetSqliteDB()

	for _, v := range tables {
		if _, err := db.Exec(v); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
