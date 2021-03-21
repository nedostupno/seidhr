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

func (s *UserService) IsHasSubsriptions(tguserID int) (bool, error) {
	return s.repo.IsHasSubsriptions(tguserID)
}
