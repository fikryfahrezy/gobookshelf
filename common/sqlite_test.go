package common

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func TestDbConnection(t *testing.T) {
	fDb := "./../data/db-test"
	sqldb, err := sql.Open("sqlite", fDb)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = sqldb.Exec(`
		drop table if exists t;
		create table t(i);
		insert into t values(42), (314);
	`); err != nil {
		t.Fatal(err)
	}

	sqldb.Close()
	_, err = os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(fDb)
	if err != nil {
		t.Fatal(err)
	}
}
