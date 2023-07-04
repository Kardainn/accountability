package server

import (
	"net/http"

	"github.com/Kardainn/accountability/backend/database"
	"github.com/gorilla/mux"
)

func (s *Server) userCreation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		database.CreateUser(s.ctx, w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func (s *Server) userIdGeneric(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	strId, ok := vars["id"]

	if !ok || strId == "" {
		w.WriteHeader(412)
		w.Write([]byte("Missing user id"))
	}

	switch r.Method {
	case "DELETE":
		database.DeleteUser(s.ctx, w, r)
		return
	case "PATCH":
		database.PatchUser(s.ctx, w, r)
		return
	case "GET":
		database.GetUser(s.ctx, w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
