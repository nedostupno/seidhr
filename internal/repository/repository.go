package repository

import "github.com/jmoiron/sqlx"

type Users interface {
	CreateUser()
	CheckUser()
	GetState()
	ChangeState()
	GetSubscriptions()
	IsHasSubsriptions()
	GetChatID()
	GetSelectedMed()
	ChangeSelectedMed()
	IsSubToThisMed()
}

type Medicoments interface {
	IsMedExist()
	GetTrueMedName()
	Subscribe()
	Unsubscribe()
	GetSubscribers()
	GetAvaliability()
	ChangeAvaliability()
	GetMedID()
	GetMedTitle()
}

type Repository struct {
	Users
	Medicoments
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		// Users:        NewUserRepo(db),
		// Medicoments: NewMedicomentsRepo(db),
	}
}
