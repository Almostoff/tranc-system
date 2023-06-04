package repository

import (
	"errors"
	"log"
	"strconv"
	"tranc/server/internal/database"
	"tranc/server/internal/entity"
	"tranc/server/internal/rabbit"
)

type Repository struct {
	db database.Database
	rb rabbit.Rabbit
}

func InitStore(db database.Database, rb rabbit.Rabbit) *Repository {
	Repository := Repository{
		db: db,
		rb: rb,
	}
	return &Repository
}

func (r Repository) CreateItem(tranc *entity.Tranc, id int) error {
	if tranc.Amount < 0 {
		balance := r.db.GetBalance(id)
		if balance < tranc.Amount*-1 {
			err := errors.New("user with id: " + strconv.Itoa(id) + " doesn`t have enough money")
			return err
		} else {
			err := r.rb.SendMessage(id, tranc)
			if err != nil {
				log.Println("can`t send msg to rabbit: ", err)
				return err
			}

			go r.updateData(id)
			return nil
		}
	} else {
		err := r.rb.SendMessage(id, tranc)
		if err != nil {
			log.Println("can`t send msg to rabbit: ", err)
			return err
		}

		go r.updateData(id)
		return nil
	}
}

func (r Repository) updateData(id int) error {
	complete := make(chan bool)
	msg, err := r.rb.ReadMessage(id, complete)
	if err != nil {
		log.Println("failed to read msg from rabbit", err)
		return err
	}
	err = r.db.CreateItem(msg, id, complete)
	if err != nil {
		log.Println("can`t update balance: ", err)
		return err
	}
	return nil
}
