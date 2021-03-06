package repository

import "github.com/jmoiron/sqlx"

type Users interface {
	Create(tguserID int, chatID int64) error
	Check(tguserID int) (bool, error)
	GetState(tguserID int) (string, error)
	ChangeState(tguserID int, state string) error
	GetChatID(tguserID int) (int, error)
	IsHasSubsriptions(tguserID int) (bool, error)
	GetSelectedMed(tguserID int) (int, error)
	ChangeSelectedMed(medicamentID, tguserID int) error
	IsSubToThisMed(tguserID int, medicamentID int) (bool, error)
	GetSubscriptions(tguserID int) ([][]string, error)
	Subscribe(tguserID int, medicamentID int) error
	Unsubscribe(tguserID int, medicamentID int) error
}

type Medicaments interface {
	IsExist(medName string) (bool, error)
	GetTrueName(medName string) (string, error)
	GetID(medTitle string) (int, error)
	GetTitle(medicamentID int) (string, error)
	InitMedList(medLines []string) error
	IsMedListExist() (bool, error)
	AreTheAnySubscriptions() (bool, error)
	GetAllMedicamentsWithSub() ([]int, error)
	GetAvailability(medicamentID int) (bool, error)
	ChangeAvailability(medicamentID int, value bool) error
	GetSubscribers(medicamentID int) ([]int, error)
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
