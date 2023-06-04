package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"tranc/server/internal/cases"
	"tranc/server/internal/entity"
)

type Controller struct {
	usecase cases.Usecase
}

func NewController(usecase cases.Usecase) *Controller {
	return &Controller{
		usecase: usecase,
	}
}

func Build(r *chi.Mux, usecase cases.Usecase) {
	ctr := NewController(usecase)

	r.Route("/tranc", func(r chi.Router) {
		r.Post("/{ID}", ctr.CreateItem)
	})

}

func (s *Controller) CreateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	uId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id: not a number"))
		return
	}
	item := &entity.Tranc{}
	err = json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = s.usecase.CreateItem(item, uId)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(err.Error()))
	}
}
