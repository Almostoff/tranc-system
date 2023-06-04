package app

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"tranc/server/config"
	"tranc/server/internal/cases"
	"tranc/server/internal/controller"
	"tranc/server/internal/database"
	"tranc/server/internal/rabbit"
	"tranc/server/internal/repository"
)

type App struct {
	cfg config.Config
}

func InitApp(cfg config.Config) *App {
	app := App{}
	app.cfg = cfg
	return &app
}

func (app *App) Run() {
	db, err := database.InitDBConn(app.cfg)
	if err != nil {
		log.Println("Error while connecting to database")
		panic(err)
	} else {
		log.Println("Connect to db successful")
	}

	rb, err := rabbit.NewRabbit()
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Connect to rMQ successful")
	}

	repository := repository.InitStore(db, *rb)
	usecase := cases.NewUseCase(repository)
	r := chi.NewRouter()

	controller.Build(r, usecase)

	err = http.ListenAndServe("localhost:8181", r)
	if err != nil {
		log.Fatal("Сервер не запустился")
		return
	}
}
