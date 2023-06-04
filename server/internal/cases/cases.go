package cases

import (
	"tranc/server/internal/entity"
	"tranc/server/internal/repository"
)

type (
	Usecase interface {
		CreateItem(*entity.Tranc, int) error
	}

	DataBase interface {
		CreateItem(tranc *entity.Tranc) error
	}
)

type usecase struct {
	rep *repository.Repository
}

func NewUseCase(rep *repository.Repository) *usecase {
	return &usecase{
		rep: rep,
	}
}

func (u *usecase) CreateItem(tranc *entity.Tranc, id int) error {
	err := u.rep.CreateItem(tranc, id)
	return err
}
