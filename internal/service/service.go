package service

import "github.com/nedostupno/seidhr/internal/repository"

type Users interface {
	Create(tguserID int, chatID int64) error
	Check(tguserID int) (bool, error)
	IsHasSubsriptions(tguserID int) (bool, error)
}

type Service struct {
	Users
}

func NewService(repo *repository.Repository) *Service {
	return &Service{Users: NewUserService(repo.Users)}
}
