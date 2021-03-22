package service

import "github.com/nedostupno/seidhr/internal/repository"

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(tguserID int, chatID int64) error {
	return s.repo.Create(tguserID, chatID)
}

func (s *UserService) Check(tguserID int) (bool, error) {
	return s.repo.Check(tguserID)
}

func (s *UserService) GetState(tguserID int) (string, error) {
	return s.repo.GetState(tguserID)
}

func (s *UserService) ChangeState(tguserID int, state string) error {
	return s.repo.ChangeState(tguserID, state)
}

func (s *UserService) IsHasSubsriptions(tguserID int) (bool, error) {
	return s.repo.IsHasSubsriptions(tguserID)
}

func (s *UserService) GetSelectedMed(tguserID int) (int, error) {
	return s.repo.GetSelectedMed(tguserID)
}

func (s *UserService) ChangeSelectedMed(medicamentID, tguserID int) error {
	return s.repo.ChangeSelectedMed(medicamentID, tguserID)
}

func (s *UserService) IsSubToThisMed(tguserID int, medicamentID int) (bool, error) {
	return s.repo.IsSubToThisMed(tguserID, medicamentID)
}

func (s *UserService) GetSubscriptions(tguserID int) ([][]string, error) {
	return s.repo.GetSubscriptions(tguserID)
}

func (s *UserService) Subscribe(tguserID int, medicamentID int) error {
	return s.repo.Subscribe(tguserID, medicamentID)
}
