package repository

import (
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// Create - создает нового польователя
func (u *UserPostgres) Create(tguserID int, chatID int64) error {
	tx := u.db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, chat_id, state, selected_med) VALUES ($1, $2, $3, $4)", tguserID, chatID, "born", 0)
	err := tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Check - проверяет наличие пользователя в базе
func (u *UserPostgres) Check(tguserID int) (bool, error) {
	var isExist bool
	err := u.db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// IsHasSubsriptions - проверяет наличие у пользователя подписок на лекарства и,
// если у него есть хоть одна подписка, возвращает true.
func (u *UserPostgres) IsHasSubsriptions(tguserID int) (bool, error) {
	var isExist bool
	err := u.db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}
