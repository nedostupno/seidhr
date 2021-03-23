package repository

import "github.com/jmoiron/sqlx"

type MedPostgres struct {
	db *sqlx.DB
}

func NewMedPostgres(db *sqlx.DB) *MedPostgres {
	return &MedPostgres{db: db}
}

// IsExist - проверяет существует ли такое лекарство в нашей базе
func (m *MedPostgres) IsExist(medName string) (bool, error) {

	var isExist bool
	err := m.db.QueryRow("SELECT exists (select 1 from medicament where title % $1)", medName).Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// GetTrueName - выводит правильное название лекарства, если пользователь
// ввел название с опечатками.
//
// Данная функция используется в связке с IsMedExist.
//
//IsMedExist - проверяет существования лекарства, а данная функция
// выдает правильное название, для дальнешей работы с Гос. Услугами
func (m *MedPostgres) GetTrueName(medName string) (string, error) {

	var trueName string
	err := m.db.QueryRow("SELECT title FROM medicament WHERE title % $1", medName).Scan(&trueName)
	if err != nil {
		return "", err
	}
	return trueName, nil
}

// GetMedID - находит id необхомодимого лекартсва
func (m *MedPostgres) GetID(medTitle string) (int, error) {
	var medicamentID int
	err := m.db.QueryRow("SELECT id FROM medicament WHERE title = $1", medTitle).Scan(&medicamentID)
	if err != nil {
		return 0, err
	}
	return medicamentID, nil
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

// GetMedTitle - находит название лекарства по его id
func (m *MedPostgres) GetTitle(medicamentID int) (string, error) {
	var medTitle string
	err := m.db.QueryRow("SELECT title FROM medicament WHERE id = $1", medicamentID).Scan(&medTitle)
	if err != nil {
		return "", err
	}
	return medTitle, nil
}

// AreTheAnySubscriptions - проверяет существование хотя бы одной подписки
// Служит для того, что бы избежать ошибок в функции CyclicMedSearch
func (m *MedPostgres) AreTheAnySubscriptions() (bool, error) {
	var isExist bool
	err := m.db.QueryRow("SELECT exists (SELECT 1 FROM subscription )").Scan(&isExist)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

// GetAllMedicamentsWithSub - Находит все лекарства, на которые подписаны пользователи
// и возварщает слайс с их id
func (m *MedPostgres) GetAllMedicamentsWithSub() ([]int, error) {
	rows, err := m.db.Query("SELECT DISTINCT medicament_id FROM subscription")
	if err != nil {
		return nil, err
	}

	subMeds := []int{}

	for rows.Next() {
		var id int

		rows.Scan(&id)
		subMeds = append(subMeds, id)
		defer rows.Close()
	}
	return subMeds, nil
}

// GetAvailability - проверяет наличие лекарства записаное в базе
func (m *MedPostgres) GetAvailability(medicamentID int) (bool, error) {
	var availible bool
	err := m.db.QueryRow("SELECT availability FROM medicament WHERE id = $1", medicamentID).Scan(&availible)
	if err != nil {
		return false, err
	}
	return availible, nil
}

// ChangeAvailability - изменяет наличие лекарства в базе
func (m *MedPostgres) ChangeAvailability(medicamentID int, value bool) error {
	tx := m.db.MustBegin()
	tx.MustExec("UPDATE medicament SET availability = $1 WHERE id = $2", value, medicamentID)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// GetSubscribers - находит пользователей подписанных на определенное лекарство,
// id которого принимается на вход, и возвращает слайс с id этих пользователей
func (m *MedPostgres) GetSubscribers(medicamentID int) ([]int, error) {
	rows, err := m.db.Query("SELECT tguser_id FROM subscription WHERE medicament_id = $1", medicamentID)
	if err != nil {
		return nil, err
	}

	users := []int{}

	for rows.Next() {
		var id int

		rows.Scan(&id)
		users = append(users, id)
		defer rows.Close()
	}
	return users, nil
}
