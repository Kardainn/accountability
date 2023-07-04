package server

import (
	"net/http"
	"strconv"

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

	id64, err := strconv.ParseUint(strId, 10, 32)
	if err != nil {
		w.WriteHeader(412)
		w.Write([]byte("User ID is cannot be parsed as uint64"))
	}
	id32 := uint32(id64)

	switch r.Method {
	case "DELETE":
		database.DeleteUser(s.ctx, w, r, id32)
		return
	case "PATCH":
		database.PatchUser(s.ctx, w, r, id32)
		return
	case "GET":
		database.GetUser(s.ctx, w, r, id32)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
