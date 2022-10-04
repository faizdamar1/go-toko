package app

import (
	"github.com/faizdamar1/go-toko/app/controllers"
	"github.com/gorilla/mux"
)

func (server *Server) inisializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
