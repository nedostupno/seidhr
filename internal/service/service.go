package service

import (
	"github.com/nedostupno/seidhr/internal/model"
	"github.com/nedostupno/seidhr/internal/repository"
)

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
	ReadFileWithMeds() ([]string, error)
	ReqMedInfo(medTitle string) (model.Jsn, error)
	ParseJSON(j model.Jsn) string
	IsErrExistInJSON(j model.Jsn) bool
	AreTheAnySubscriptions() (bool, error)
	GetAllMedicamentsWithSub() ([]int, error)
	GetAvailability(medicamentID int) (bool, error)
	ChangeAvailability(medicamentID int, value bool) error
	GetSubscribers(medicamentID int) ([]int, error)
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
