package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/speedmancs/vmmanager/middleware"
	"github.com/speedmancs/vmmanager/model"
	"github.com/speedmancs/vmmanager/util"
)

type App struct {
	Router *mux.Router
}

func (app *App) getVMs(w http.ResponseWriter, r *http.Request) {
	vms, err := model.GetVMs()
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondWithJSON(w, http.StatusOK, vms)
	}
}

func (app *App) getVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	vm, err := model.GetVM(id)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) updateVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	vm, err := model.UpdateVM(id, r.Body)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
	token, err := middleware.Login(r.Body)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, err.Error())
	} else {
		tokenObj := middleware.Token{
			Token: token,
		}
		util.RespondWithJSON(w, http.StatusOK, tokenObj)
	}
}

func (app *App) registerVM(w http.ResponseWriter, r *http.Request) {
	vm, err := model.RegisterVM(r.Body)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondWithJSON(w, http.StatusOK, vm)
	}
}

func (app *App) deleteVM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Invalid ID %s", vars["id"]))
		return
	}

	err = model.DeleteVM(id)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		util.RespondWithInfo(w, http.StatusOK, fmt.Sprintf("VM %s deleted", vars["id"]))
	}
}

func (app *App) Initialize(logFile string) {
	app.Router = mux.NewRouter()
	app.Router.Use(middleware.RequestIDMiddleware, middleware.LoggingMiddleware)
	app.Router.HandleFunc("/login", app.login).Methods("POST")

	subRouter := app.Router.PathPrefix("/vm").Subrouter()
	subRouter.Use(middleware.AuthMiddleware)
	subRouter.HandleFunc("/", app.registerVM).Methods("POST")
	subRouter.HandleFunc("/all", app.getVMs).Methods("GET")
	subRouter.HandleFunc("/{id:[0-9]+}", app.getVM).Methods("GET")
	subRouter.HandleFunc("/{id:[0-9]+}", app.updateVM).Methods("PUT")
	subRouter.HandleFunc("/{id:[0-9]+}", app.deleteVM).Methods("DELETE")

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func (app *App) Run(port string) {
	log.Println("Starting service...")
	log.Fatal(http.ListenAndServe(port, app.Router))
}
