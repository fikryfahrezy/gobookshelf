package db

import "testing"

func TestDbConnection(t *testing.T) {
	fd := "./../data/db-test"

	db, err := InitSqliteTestDB(fd)
	if err != nil {
		t.FailNow()
	}

	if _, err = db.Exec(`
		drop table if exists t;
		create table t(i);
		insert into t values(42), (314);
	`); err != nil {
		t.FailNow()
	}

	if err = RemoveSqliteTestDB(db, fd); err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}
