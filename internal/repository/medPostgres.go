package repository

import "github.com/jmoiron/sqlx"

type MedPostgres struct {
	db *sqlx.DB
}

func NewMedPostgres(db *sqlx.DB) *MedPostgres {
	return &MedPostgres{db: db}
}

// InitMedList - инициализирует список льготных лекарств в базе данных
func (m *MedPostgres) InitMedList(medLines []string) error {
	tx := m.db.MustBegin()

	for _, med := range medLines {
		tx.MustExec("INSERT INTO medicament (title, availability) VALUES ($1, $2)", med, false)
	}
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// IsMedListExist - проверяет заполнена ли таблица medicament.
// Служит для того что бы не пытаться каждый раз заполнять бд значениями из файла drugs.txt
func (m *MedPostgres) IsMedListExist() (bool, error) {
	var isExist bool
	err := m.db.QueryRow("SELECT exists (select 1 from medicament)").Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}
