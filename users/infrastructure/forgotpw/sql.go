package forgotpw

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/users/domain/users"
)

type ForgotPassRepository struct {
	Db *sql.DB
}

func (r ForgotPassRepository) Insert(fp users.ForgotPassModel) {
	fp.Id = common.RandString(4)

	q := `
		INSERT INTO user_forgot_pass(id, email, code, is_claimed) VALUES( ?, ?, ?, ? )
	`

	_, err := r.Db.Exec(q, fp.Id, fp.Email, fp.Code, fp.IsClaimed)
	if err != nil {
		fmt.Println(err)
	}
}

func (r ForgotPassRepository) ReadByCode(k string) (users.ForgotPassModel, error) {
	var fp users.ForgotPassModel
	q := `
		SELECT id, email, code, is_claimed FROM user_forgot_pass WHERE code=?
	`
	err := r.Db.QueryRow(q, k).Scan(&fp.Id, &fp.Email, &fp.Code, &fp.IsClaimed)

	switch {
	case err == sql.ErrNoRows:
		return fp, errors.New("forgot pass not found")
	case err != nil:
		return fp, err
	default:
		return fp, nil
	}
}

func (r ForgotPassRepository) Update(fp users.ForgotPassModel) (users.ForgotPassModel, error) {
	result, err := r.Db.Exec("UPDATE user_forgot_pass SET is_claimed = ? WHERE id = ? AND is_claimed=0", fp.IsClaimed, fp.Id)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fp, err
	}

	if rows != 1 {
		return fp, errors.New("forgot pass not found")
	}

	return fp, nil
}
