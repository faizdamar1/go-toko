package app

import (
	"net/http"

	"github.com/faizdamar1/go-toko/app/controllers"
	"github.com/gorilla/mux"
)

func (server *Server) inisializeRoutes() {
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/", controllers.Home).Methods("GET")

	staticFileDirectory := http.Dir("./assets")
	staticFileHandler := http.StripPrefix("/public", http.FileServer(staticFileDirectory))

	server.Router.PathPrefix("/public").Handler(staticFileHandler).Methods("GET")
}
