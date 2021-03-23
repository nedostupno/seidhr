package repository

import (
	"strconv"

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

// GetState - проверяет состояние пользователя
func (u *UserPostgres) GetState(tguserID int) (string, error) {
	var state string
	err := u.db.QueryRow("SELECT state from tguser where id=$1", tguserID).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

// ChangeState - изменяет состояние пользователя
func (u *UserPostgres) ChangeState(tguserID int, state string) error {
	tx := u.db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, tguserID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetChatID - находит пользователя и возвращает его chatID
func (u *UserPostgres) GetChatID(tguserID int) (int, error) {
	var chatID int
	err := u.db.QueryRow("SELECT chat_id FROM tguser WHERE id = $1", tguserID).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}

// GetSelectedMed - получает id лекарства выбранного пользователем в данный момент
func (u *UserPostgres) GetSelectedMed(tguserID int) (int, error) {
	var medicamentID int
	err := u.db.QueryRow("SELECT selected_med FROM tguser WHERE id = $1", tguserID).Scan(&medicamentID)
	if err != nil {
		return 0, err
	}
	return medicamentID, nil
}

// ChangeSelectedMed - меняет выбранное пользователем лекарство
func (u *UserPostgres) ChangeSelectedMed(medicamentID, tguserID int) error {
	tx := u.db.MustBegin()
	tx.MustExec("UPDATE tguser SET selected_med = $1 WHERE id = $2", medicamentID, tguserID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// IsSubToThisMed - проверяет подписан ли пользователь на данное лекарство
func (u *UserPostgres) IsSubToThisMed(tguserID int, medicamentID int) (bool, error) {
	var isExist bool
	err := u.db.QueryRow("SELECT exists (SELECT 1 FROM subscription WHERE tguser_id=$1 AND medicament_id=$2)", tguserID, medicamentID).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// GetSubscriptions - находит все подписки пользователя и возвращает [][]string, где
// [[id title] [id title] [id title]]
func (u *UserPostgres) GetSubscriptions(tguserID int) ([][]string, error) {
	rows, err := u.db.Query("SELECT id, title from medicament INNER JOIN subscription on medicament.id=subscription.medicament_id WHERE subscription.tguser_id = $1", tguserID)
	if err != nil {
		return nil, err
	}

	subscriptions := [][]string{}

	for rows.Next() {
		var id int
		var title string

		rows.Scan(&id, &title)
		subscriptions = append(subscriptions, []string{strconv.Itoa(id), title})
		defer rows.Close()
	}
	return subscriptions, nil
}

func (u *UserPostgres) Subscribe(tguserID int, medicamentID int) error {
	tx := u.db.MustBegin()
	tx.MustExec("INSERT INTO subscription (tguser_id, medicament_id) VALUES ($1, $2)", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Unsubscribe - отменяет у пользьователя подписку на лекарство
func (u *UserPostgres) Unsubscribe(tguserID int, medicamentID int) error {
	tx := u.db.MustBegin()
	tx.MustExec("DELETE FROM subscription WHERE tguser_id = $1 AND medicament_id = $2", tguserID, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
