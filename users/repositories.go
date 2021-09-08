package users

import (
	"database/sql"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/fikryfahrezy/gobookshelf/common"
	"github.com/fikryfahrezy/gobookshelf/db"
)

type userModel struct {
	Id       string
	Email    string
	Password string
	Name     string
	Region   string
	Street   string
}

func (um *userModel) Save() (userModel, error) {
	um.Id = common.RandString(7)
	ur := users.Insert(*um)

	return *um, ur
}

func (um *userModel) Update(nu userModel) (userModel, error) {
	if nu.Email != "" {
		um.Email = nu.Email
	}

	if nu.Password != "" {
		um.Password = nu.Password
	}

	if nu.Name != "" {
		um.Name = nu.Name
	}

	if nu.Region != "" {
		um.Region = nu.Region
	}

	if nu.Street != "" {
		um.Street = nu.Street
	}

	nu, ok := users.Update(*um)

	return nu, ok
}

type userDB struct {
	users map[time.Time]userModel
	lock  sync.RWMutex
}

func (udb *userDB) Insert(u userModel) error {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	for _, v := range udb.users {
		if v.Email == u.Email {
			return errors.New("email already exist")
		}
	}

	udb.users[time.Now()] = u

	return nil
}

func (udb *userDB) ReadByEmail(k string) (userModel, error) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.users {
		if v.Email == k {
			return v, nil
		}
	}

	return userModel{}, errors.New("user not found")
}

func (udb *userDB) ReadById(k string) (userModel, error) {
	udb.lock.RLock()
	defer udb.lock.RUnlock()

	for _, v := range udb.users {
		if v.Id == k {
			return v, nil
		}
	}

	return userModel{}, errors.New("user not found")
}

func (udb *userDB) Update(u userModel) (userModel, error) {
	udb.lock.Lock()
	defer udb.lock.Unlock()

	for _, v := range udb.users {
		if v.Email == u.Email && v.Id != u.Id {
			return userModel{}, errors.New("user not found")
		}
	}

	for i, v := range udb.users {
		if v.Id == u.Id {
			udb.users[i] = u
			return udb.users[i], nil

		}
	}

	return userModel{}, errors.New("user not found")
}

var users = userDB{users: make(map[time.Time]userModel)}

type forgotPassModel struct {
	Id        string
	Email     string
	Code      string
	IsClaimed bool
}

func (fpM *forgotPassModel) Save() {
	fpM.Id = common.RandString(4)

	ForgotPasses.Insert(*fpM)
}

func (fpM *forgotPassModel) Update(nfpM forgotPassModel) (forgotPassModel, error) {
	if nfpM.Code != "" {
		fpM.Code = nfpM.Code
	}

	if nfpM.Email != "" {
		fpM.Email = nfpM.Email
	}

	if nfpM.IsClaimed != fpM.IsClaimed {
		fpM.IsClaimed = nfpM.IsClaimed
	}

	nfpM, ok := ForgotPasses.Update(*fpM)

	return nfpM, ok
}

type forgotPassDB struct{}

func (fpdb *forgotPassDB) Insert(fp forgotPassModel) {
	sqliteDB := db.GetSqliteDB()
	q := `
		INSERT INTO user_forgot_pass(id, email, code, is_claimed) VALUES( ?, ?, ?, ? )
	`

	_, err := sqliteDB.Exec(q, fp.Id, fp.Email, fp.Code, fp.IsClaimed)
	if err != nil {
		log.Fatal(err)
	}
}

func (fpdb *forgotPassDB) ReadByCode(k string) (forgotPassModel, error) {
	var fp forgotPassModel

	sqliteDB := db.GetSqliteDB()
	q := `
		SELECT id, email, code, is_claimed FROM user_forgot_pass WHERE code=?
	`
	err := sqliteDB.QueryRow(q, k).Scan(&fp.Id, &fp.Email, &fp.Code, &fp.IsClaimed)

	switch {
	case err == sql.ErrNoRows:
		return fp, errors.New("forgot pass not found")
	case err != nil:
		return fp, err
	default:
		return fp, nil
	}
}

func (fpdb *forgotPassDB) Update(fp forgotPassModel) (forgotPassModel, error) {
	sqliteDB := db.GetSqliteDB()
	result, err := sqliteDB.Exec("UPDATE user_forgot_pass SET is_claimed = ? WHERE id = ? AND is_claimed=0", fp.IsClaimed, fp.Id)
	if err != nil {
		log.Fatal(err)
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

var ForgotPasses = forgotPassDB{}
