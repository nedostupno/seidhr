package repository

import "github.com/jmoiron/sqlx"

type Users interface {
	Create(tguserID int, chatID int64) error
	Check(tguserID int) (bool, error)
	GetState(tguserID int) (string, error)
	ChangeState(tguserID int, state string) error
	IsHasSubsriptions(tguserID int) (bool, error)
	GetSelectedMed(tguserID int) (int, error)
	ChangeSelectedMed(medicamentID, tguserID int) error
	IsSubToThisMed(tguserID int, medicamentID int) (bool, error)
	GetSubscriptions(tguserID int) ([][]string, error)
}

type Medicaments interface {
	IsExist(medName string) (bool, error)
	GetTrueName(medName string) (string, error)
	GetID(medTitle string) (int, error)
	InitMedList(medLines []string) error
	IsMedListExist() (bool, error)
}

type Repository struct {
	Users
	Medicaments
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users:       NewUserPostgres(db),
		Medicaments: NewMedPostgres(db),
	}
}
