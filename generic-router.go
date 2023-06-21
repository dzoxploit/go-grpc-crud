package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzoxploit/go-grpc-crud/repos"

	"github.com/gorilla/mux"
)

type GenericRouter[TT any, T repos.GenericRepo[TT]] struct {
	muxBase string
	repo    *repos.GenericRepo[TT]
}

func (rtr *GenericRouter[TT, T]) handle(w http.ResponseWriter, r *http.Request) {
	idLong := mux.Vars(r)["id"]
	if r.URL.EscapedPath() != rtr.muxBase {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(idLong) != 0 {
			item, err := (*(rtr.repo)).GetOne(string(idLong))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode("Entity not found!")
				return
			}
			json.NewEncoder(w).Encode(&item)
		} else {
			items := (*(rtr.repo)).GetList()
			json.NewEncoder(w).Encode(&items)
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		var model TT
		json.NewDecoder(r.Body).Decode(&model)
		(*(rtr.repo)).Create(model)
		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		var model TT
		json.NewDecoder(r.Body).Decode(&model)
		_, err := (*(rtr.repo)).Update(string(idLong), model)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Entity not found!")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if len(idLong) != 0 {
			_, err := (*(rtr.repo)).DeleteOne(string(idLong))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode("Entity not found!")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (rtr *GenericRouter[TT, T]) registerRoutes(mux *mux.Router) {
	mux.HandleFunc(rtr.muxBase, rtr.handle)
	mux.HandleFunc(fmt.Sprintf("%v/{id}", rtr.muxBase), rtr.handle)
}

func NewGenericRouter[TT any, T repos.GenericRepo[TT]](muxBase string, mux *mux.Router, repo *repos.GenericRepo[TT]) *GenericRouter[TT, T] {
	router := GenericRouter[TT, T]{
		muxBase: muxBase,
		repo:    repo,
	}
	router.registerRoutes(mux)
	return &router
}