package common

func MigrateSqliteDB() string {
	return `
			CREATE TABLE IF NOT EXISTS user_forgot_pass (
				id TEXT,
				email TEXT,
				code TEXT,
				is_claimed INTEGER
			);
		`
}
