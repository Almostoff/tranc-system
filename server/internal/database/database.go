package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
	"tranc/server/config"
	"tranc/server/internal/entity"
)

type Database struct {
	db *sqlx.DB
}

func InitDBConn(cfg config.Config) (Database, error) {
	dbConn := Database{}
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Addr, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBname)
	dbConn.db, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return Database{}, err
	}

	return dbConn, err
}

func (d Database) GetBalance(id int) (res int) {
	err := d.db.Get(&res, `select (select coalesce(sum(balance), 0) from clients where id=$1)`, id)
	if err != nil {
		panic(err)
	}
	return res
}

func (d Database) CreateItem(data chan *entity.Tranc, id int, complete chan bool) error {
	for v := range data {
		if v.Amount < 0 {
			balance := d.GetBalance(id)
			if balance < v.Amount*-1 {
				return errors.New("not enough money")
			}
		}
		q := `INSERT INTO invoice (u_id, amount, created_at)
          VALUES ($1, $2, $3)`
		_, err := d.db.Exec(q, id, v.Amount, time.Now())
		if err != nil {
			return err
		}
		complete <- true
		q1 := `UPDATE clients SET balance = (SELECT SUM(amount) FROM invoice WHERE u_id = $1) WHERE id = $1
	  RETURNING id`

		_, err = d.db.Exec(q1, id)
		if err != nil {
			return err
		}
	}
	return nil
}
