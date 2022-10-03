package app

import (
	"github.com/faizdamar1/go-toko/app/controllers"
)

func (server *Server) inisializeRoutes() {
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")
}
