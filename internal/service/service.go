package service

import "github.com/nedostupno/seidhr/internal/repository"

type Users interface {
	Create(tguserID int, chatID int64) error
	Check(tguserID int) (bool, error)
	IsHasSubsriptions(tguserID int) (bool, error)
}

type Medicaments interface {
	InitMedList(medLines []string) error
	IsMedListExist() (bool, error)
	ReadFileWithMeds() ([]string, error)
}

type Service struct {
	Users
	Medicaments
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Users:       NewUserService(repo.Users),
		Medicaments: NewMedService(repo.Medicaments),
	}
}
