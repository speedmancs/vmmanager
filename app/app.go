package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/speedmancs/vmmanager/model"
)

type App struct {
	Router *mux.Router
}

func (app *App) getVMs(w http.ResponseWriter, r *http.Request) {
	model.GetVMs()
}

func (app *App) getVM(w http.ResponseWriter, r *http.Request) {

}

func (app *App) updateVM(w http.ResponseWriter, r *http.Request) {

}

func (app *App) registerVM(w http.ResponseWriter, r *http.Request) {

}

func (app *App) deleteVM(w http.ResponseWriter, r *http.Request) {

}

func (app *App) Initialize() {
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/vms", app.getVMs).Methods("GET")
	app.Router.HandleFunc("/vm", app.registerVM).Methods("POST")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.getVM).Methods("GET")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.updateVM).Methods("PUT")
	app.Router.HandleFunc("/vm/{id:[0-9]+}", app.deleteVM).Methods("DELETE")
}

func (app *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, app.Router))
}
