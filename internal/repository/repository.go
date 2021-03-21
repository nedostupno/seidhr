package repository

import "github.com/jmoiron/sqlx"

type Users interface {
	Create(tguserID int, chatID int64) error
	Check(tguserID int) (bool, error)
	IsHasSubsriptions(tguserID int) (bool, error)
}

// type Medicaments interface {
// }

type Repository struct {
	Users
	// Medicaments
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUserPostgres(db),
		// Medicoments: NewMedicaments(db),
	}
}
